package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/dave/gopackages"
	"github.com/dave/jennifer/jen"
	"github.com/dave/rebecca"
)

var flags struct {
	pkg, input, output, literals string
}

func init() {
	flag.StringVar(&flags.pkg, "package", "", "Package to scan")
	flag.StringVar(&flags.input, "input", "README.md.tpl", "Input file")
	flag.StringVar(&flags.output, "output", "", "Output file, defaults to the input without the .tpl suffix")
	flag.StringVar(&flags.literals, "literals", "", "Output Go file, containing map of doc literals")
}

func abort(s string, vv ...interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: "+s, vv...)
	os.Exit(2)
}

func main() {
	flag.Parse()

	if flags.input == "" {
		flag.PrintDefaults()
		return
	}
	if flags.output == "" {
		flags.output = strings.TrimSuffix(flags.input, ".tpl")
	}
	if flags.output == flags.input {
		abort("input and output both point at the same file, %s\n", flags.output)
		return
	}

	if flags.pkg == "" {
		wd, err := os.Getwd()
		if err != nil {
			abort("can't auto-detect package, %s\n", err.Error())
			return
		}
		flags.pkg, err = gopackages.GetPackageFromDir(os.Getenv("GOPATH"), wd)
		if err != nil {
			abort("can't auto-detect package, %s\n", err.Error())
			return
		} else if flags.pkg == "" {
			abort("can't find package at current dir (%s) and no package specified with 'package' flag.\n", wd)
			return
		}
	}

	dir, err := gopackages.GetDirFromPackage(os.Environ(), os.Getenv("GOPATH"), flags.pkg)
	if err != nil {
		abort("can't parse package directory, %s\n", err.Error())
		return
	}

	m, err := rebecca.NewCodeMap(flags.pkg, dir)
	if err != nil {
		abort("can't init code map, %s\n", err.Error())
		return
	}

	funcMap := template.FuncMap{
		"example": m.ExampleFunc(false),
		"code":    m.ExampleFunc(true),
		"output":  m.OutputFunc,
		"doc":     m.DocFunc,
	}

	tpl, err := template.New("main").Funcs(funcMap).ParseFiles(flags.input)
	if err != nil {
		abort("can't parse template, %s\n", err.Error())
		return
	}

	buf := &bytes.Buffer{}
	if err := tpl.ExecuteTemplate(buf, flags.input, nil); err != nil {
		abort("can't process template, %s\n", err.Error())
		return
	}
	if err := ioutil.WriteFile(flags.output, buf.Bytes(), 0644); err != nil {
		abort("can't write output, %s\n", err.Error())
		return
	}

	if flags.literals != "" {
		f := jen.NewFile(m.Name)
		f.Var().Id("doc").Op("=").Map(jen.String()).String().Values(
			jen.DictFunc(func(d jen.Dict) {
				for k, v := range m.Comments {
					d[jen.Lit(k)] = jen.Lit(strings.TrimSpace(v))
				}
			}),
		)
		if err := f.Save(flags.literals); err != nil {
			abort("can't write literals file, %s\n", err.Error())
			return
		}
	}

}
