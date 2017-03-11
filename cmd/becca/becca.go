package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"text/template"

	"log"

	"github.com/davelondon/gopackages"
	"github.com/davelondon/rebecca"
)

func main() {

	pkgFlag := flag.String("package", "", "Package to scan")
	flag.Parse()
	pkg := *pkgFlag

	if pkg == "" {
		wd, _ := os.Getwd()
		pkg, _ := gopackages.GetPackageFromDir(os.Getenv("GOPATH"), wd)
		if pkg == "" {
			log.Fatalf("Can't find package at current dir (%s) and no package specified with 'package' flag.", wd)
		}
	}

	dir, err := gopackages.GetDirFromPackage(os.Environ(), os.Getenv("GOPATH"), pkg)
	if err != nil {
		log.Fatal(err)
	}

	m, err := rebecca.NewCodeMap(pkg, dir)
	if err != nil {
		log.Fatal(err)
	}

	funcMap := template.FuncMap{
		"example": m.ExampleFunc(false),
		"code":    m.ExampleFunc(true),
		"output":  m.OutputFunc,
		"doc":     m.DocFunc,
	}

	tpl := template.Must(template.New("main").Funcs(funcMap).ParseGlob("README.md.tpl"))

	buf := &bytes.Buffer{}
	if err := tpl.ExecuteTemplate(buf, "README.md.tpl", nil); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("README.md", buf.Bytes(), 0644); err != nil {
		log.Fatal(err)
	}
}
