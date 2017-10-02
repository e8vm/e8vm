[![BuildStatus](https://travis-ci.org/shanhuio/smlvm.png?branch=master)](https://travis-ci.org/shanhuio/smlvm)

# Small VM

Small Virtual Machine (`smlvm`) provides [a simple programming language][1]
that compiles to a simulated, simple virtual machine. It is
essentially a subset of Go programming language.

This repository does not use any third party libraries; it depends on
only Go standard library. The compiler is all written from scratch,
and does not use LLVM.

[Try in playground](https://g.smallrepo.com/play)

Shanhuio provides both a [cloud IDE](https://g.smallrepo.com/),
and local environment together with [Small Home](https://github.com/shanhuio/smlhome) for developement.

[G language introduction][1]

[1]: https://github.com/shanhuio/smlvm/wiki/G-Language-Introduction

## To Use

### Install

```
go get -u shanhu.io/smlvm/...
```

### make

The project comes with a `makefile`, which formats the code files,
check lints, check circular dependencies and build tags. Running the
`makefile` requires installing some tools.

```
go get -u shanhu.io/tools/...
go get -u github.com/golang/lint/golint
go get -u github.com/jstemmer/gotags
```

## What is Small VM

### Approachable and Evolving Compiler

This project is written in Go. Each source file in this project has no
more than 300 lines (80 max per line). Also there are no circular
dependencies among files. As a result, the project architecture can be
[visualized](https://shanhu.io/smlvm).

[Package Docs](https://godoc.org/shanhu.io/smlvm).

We hope that our design will make it easier for people to understand
and add new features to the compiler and make it better.

### The Language Targets Comprehension

Similar to Small VM, In G language, we have rules to keep the code
clean: no circular dependency among files, no more then 300 lines each
file, no more than 80 characters each.  For example, the architect of
the std G language can be found [here](https://g.smallrepo.com/r/std)
Together with the simple syntax system, we want create a language that
targets code comprehension.  We believe that readable code is
changeable code, and can continuously evolve.  We are also creating a
[cloud IDE](https://g.smallrepo.com/) for G language users to share
and read each others code.  Once code can be easily understood, online
IDE's can form a community with network effects, and developers can
easily customize a code -- their own or not -- to handle their special
needs.

## Copyright and License

Copyright by Shanhu contributors. Licence Apache.
