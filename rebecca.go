// Package rebecca is a readme generator.
package rebecca

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"regexp"
	"strconv"
	"strings"
)

func NewCodeMap(pkg string, dir string) (*CodeMap, error) {
	m := &CodeMap{
		pkg:      pkg,
		dir:      dir,
		Examples: map[string]*doc.Example{},
		Comments: map[string]string{},
	}
	if err := m.scanDir(); err != nil {
		return nil, err
	}
	return m, nil
}

type CodeMap struct {
	pkg      string
	dir      string
	fset     *token.FileSet
	Examples map[string]*doc.Example
	Comments map[string]string
}

func (m *CodeMap) ExampleFunc(plain bool) func(in string) string {
	return func(in string) string {
		e, ok := m.Examples[in]
		if !ok {
			panic(fmt.Sprintf("Example %s not found.", in))
		}
		buf := &bytes.Buffer{}
		if plain {
			printer.Fprint(buf, m.fset, e.Code)
			out := buf.String()
			if strings.HasSuffix(out, "\n\n}") {
				// fix annoying line-feed before end brace
				out = out[:len(out)-2] + "}"
			}
			return out
		}
		if bs, ok := e.Code.(*ast.BlockStmt); ok {
			for _, s := range bs.List {
				printer.Fprint(buf, m.fset, s)
				buf.WriteString("\n")
			}
		} else {
			printer.Fprint(buf, m.fset, e.Code)
		}
		quotes := "```"
		return fmt.Sprintf(`%sgo
%s
// Output:
// %s
%s`,
			quotes,
			strings.Trim(buf.String(), "\n"),
			strings.Replace(strings.Trim(e.Output, "\n"), "\n", "\n// ", -1),
			quotes)
	}
}

func (m *CodeMap) OutputFunc(in string) string {
	e, ok := m.Examples[in]
	if !ok {
		panic(fmt.Sprintf("Example %s not found.", in))
	}
	return strings.Trim(e.Output, "\n")
}

var docRegex = regexp.MustCompile(`(\w+)\[([0-9:, ]+)\]`)

func (m *CodeMap) DocFunc(in string) string {

	if matches := docRegex.FindStringSubmatch(in); matches != nil {
		id := matches[1]
		c, ok := m.Comments[id]
		if !ok {
			panic(fmt.Sprintf("Doc for %s not found in %s.", id, in))
		}
		return extractSections(in, matches[2], c)
	}

	c, ok := m.Comments[in]
	if !ok {
		panic(fmt.Sprintf("Doc for %s not found.", in))
	}
	return strings.Trim(c, "\n")
}

var bothRegex = regexp.MustCompile(`^(\d+):(\d+)$`)
var fromRegex = regexp.MustCompile(`^(\d+):$`)
var toRegex = regexp.MustCompile(`^:(\d+)$`)
var singleRegex = regexp.MustCompile(`^(\d+)$`)

func mustInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		// shoulnd't get here because the string has passed a regex
		panic(err)
	}
	return i
}

func checkBounds(start, end, length int, spec string) {

	if end == 0 {
		panic(fmt.Sprintf("End must be greater than 0 in %s", spec))
	}

	if start >= length {
		panic(fmt.Sprintf("Index %d out of range (length %d) in %s", start, length, spec))
	}

	if end >= length {
		panic(fmt.Sprintf("Index %d out of range (length %d) in %s", end, length, spec))
	}

	if end > -1 && start >= end {
		panic(fmt.Sprintf("Start must be less than end in %s", spec))
	}
}

func extractSections(full string, sections string, comment string) string {

	var sentances []string
	for _, s := range strings.Split(comment, ".") {
		// ignore empty sentances
		trimmed := strings.Trim(s, " \n")
		if trimmed != "" {
			sentances = append(sentances, s)
		}
	}

	var out string
	for _, section := range strings.Split(sections, ",") {
		var arr []string
		if matches := bothRegex.FindStringSubmatch(section); matches != nil {
			// "i:j"
			checkBounds(mustInt(matches[1]), mustInt(matches[2]), len(sentances), full)
			arr = sentances[mustInt(matches[1]):mustInt(matches[2])]
		} else if matches := fromRegex.FindStringSubmatch(section); matches != nil {
			// "i:"
			checkBounds(mustInt(matches[1]), -1, len(sentances), full)
			arr = sentances[mustInt(matches[1]):]
		} else if matches := toRegex.FindStringSubmatch(section); matches != nil {
			// ":i"
			checkBounds(-1, mustInt(matches[1]), len(sentances), full)
			arr = sentances[:mustInt(matches[1])]
		} else if matches := singleRegex.FindStringSubmatch(section); matches != nil {
			// "i"
			checkBounds(mustInt(matches[1]), -1, len(sentances), full)
			arr = []string{sentances[mustInt(matches[1])]}
		} else {
			panic(fmt.Sprintf("Invalid section %s in %s", section, full))
		}

		for _, s := range arr {
			s1 := strings.Trim(s, " \n")
			if s1 != "" {
				out += s + "."
			}
		}
	}
	return strings.Trim(out, " ")
}

func (m *CodeMap) scanTests(name string, p *ast.Package) error {
	for _, f := range p.Files {
		examples := doc.Examples(f)
		for _, ex := range examples {
			m.Examples["Example"+ex.Name] = ex
		}
	}
	return nil
}

func (m *CodeMap) scanPkg(name string, p *ast.Package) error {
	for _, f := range p.Files {
		for _, d := range f.Decls {
			switch d := d.(type) {
			case *ast.FuncDecl:
				if d.Doc.Text() == "" {
					continue
				}
				if d.Recv == nil {
					// function
					//fmt.Println(d.Name, d.Doc.Text())
					name := fmt.Sprint(d.Name)
					m.Comments[name] = d.Doc.Text()
				} else {
					// method
					e := d.Recv.List[0].Type
					if se, ok := e.(*ast.StarExpr); ok {
						// if the method receiver has a *, discard it.
						e = se.X
					}
					b := &bytes.Buffer{}
					printer.Fprint(b, m.fset, e)
					//fmt.Printf("%s.%s %s", b.String(), d.Name, d.Doc.Text())
					name := fmt.Sprintf("%s.%s", b.String(), d.Name)
					m.Comments[name] = d.Doc.Text()
				}
			case *ast.GenDecl:
				if d.Doc.Text() == "" {
					continue
				}
				switch s := d.Specs[0].(type) {
				case *ast.TypeSpec:
					//fmt.Println(s.Name, d.Doc.Text())
					name := fmt.Sprint(s.Name)
					m.Comments[name] = d.Doc.Text()
					if t, ok := s.Type.(*ast.StructType); ok {
						for _, f := range t.Fields.List {
							if f.Doc.Text() == "" {
								continue
							}
							if f.Names[0].IsExported() {
								fieldName := fmt.Sprint(name, ".", f.Names[0])
								m.Comments[fieldName] = f.Doc.Text()
							}
						}
					}
				case *ast.ValueSpec:
					//fmt.Println(s.Names[0], d.Doc.Text())
					if len(s.Names) == 0 {
						continue
					}
					name := fmt.Sprint(s.Names[0])
					m.Comments[name] = d.Doc.Text()
				}
			}
		}
	}
	return nil
}

func (m *CodeMap) scanDir() error {
	// Create the AST by parsing src.
	m.fset = token.NewFileSet() // positions are relative to fset
	pkgs, err := parser.ParseDir(m.fset, m.dir, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	for name, p := range pkgs {
		if strings.HasSuffix(name, "_test") {
			if err := m.scanTests(name, p); err != nil {
				return err
			}
		}
		if err := m.scanPkg(name, p); err != nil {
			return err
		}
	}

	return nil
}
