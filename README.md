# Rebecca

Rebecca is a readme generator

# Install

```
go get -u github.com/davelondon/rebecca/cmd/becca
```

# Usage

```
becca [-package={your-package}]
```

Rebecca will parse `README.md.tpl` and overwrite `README.md` with the result. See 
[README.md.tpl](https://github.com/davelondon/jennifer/blob/master/README.md.tpl) 
in the [jennifer](https://github.com/davelondon/jennifer) project for examples.
 
The package specified on the command line is parsed and examples and 
documentation comments are extracted. If no package is specified, it is 
detected from the current working directory. Rebecca uses the Go template 
library, and adds some custom template functions:  

# Example, Code, Output

```
{{ "ExampleFoo" | example }}
```

This prints the code and expected output for the `ExampleFoo` example.
  
# Doc

```
{{ "Foo" | doc }}
```

This prints the documentation comment for `Foo`. All package level declarations 
are supported (func, var, const etc.)

```
{{ "Foo.Bar" | doc }}
```

This prints the documentation comment for the `Bar` member of the `Foo` type. 
Methods and struct fields are supported.

You can also specify which sentances to print, using Go slice notation:

```
{{ "Foo[i]" | doc }}
{{ "Foo[i:j]" | doc }}
{{ "Foo[i:]" | doc }}
{{ "Foo[:i]" | doc }}
```

# Code, Output

```
{{ "ExampleFoo" | code }}
```

This prints just the code for the `ExampleFoo` example.

```
{{ "ExampleFoo" | output }}
```

This prints just the expected output for the `ExampleFoo` example.