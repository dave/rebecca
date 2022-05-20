# Rebecca

Rebecca is a readme generator. 

Managing the GitHub readme for your Go project can be a drag. When creating 
[jennifer](https://github.com/dave/jennifer) I found examples copied into 
the markdown would get out of date, and documentation was duplicated. I created 
rebecca to solve this: see 
[README.md.tpl](https://github.com/dave/jennifer/blob/master/README.md.tpl) 
in the [jennifer](https://github.com/dave/jennifer) repo for a real world 
example.

# Install

```
go install github.com/dave/rebecca/cmd/becca@latest
```

# Usage

```
becca [-package={your-package}]
```

Rebecca will read `README.md.tpl` and overwrite `README.md` with the rendered 
template. The package specified on the command line is parsed (if no package is 
specified, it is detected from the current working directory). 

The package is scanned for examples and documentation. Rebecca uses the Go 
template library, and adds some custom template functions:  

# Example

```
{{ "ExampleFoo" | example }}
```

This prints the code and expected output for the `ExampleFoo` example.
  
# Doc

```
{{ "Foo" | doc }}
```

This prints the documentation for `Foo`. All package level declarations are 
supported (`func`, `var`, `const` etc.)

```
{{ "Foo.Bar" | doc }}
```

This prints the documentation for the `Bar` member of the `Foo` type. Methods 
and struct fields are supported.

You can also specify which sentances to print, using Go slice notation:

```
{{ "Foo[i]" | doc }}
{{ "Foo[i:j]" | doc }}
{{ "Foo[i:]" | doc }}
{{ "Foo[:i]" | doc }}
```

See [here](https://github.com/dave/jennifer/blob/5f1e5084f7fff920e11d5b9098e5ae8089136a1a/README.md.tpl#L51-L58) and [here](https://github.com/dave/jennifer/blob/5f1e5084f7fff920e11d5b9098e5ae8089136a1a/README.md.tpl#L286-L299) for real-world examples of this.

# Code, Output

```
{{ "ExampleFoo" | code }}
```

This prints just the code for the `ExampleFoo` example.

```
{{ "ExampleFoo" | output }}
```

This prints just the expected output for the `ExampleFoo` example.
