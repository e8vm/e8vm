[![BuildStatus](https://travis-ci.org/shanhuio/smlvm.png?branch=master)](https://travis-ci.org/shanhuio/smlvm)

# Small VM

```
go get -u shanhu.io/smlvm/...
```

(Note that we use a custom domain rather than `github.com`.)

`smlvm` is a small, simulated virtual machine with a 32-bit
instruction set. It also comes with
[a simple programming langauge called G][1]. The language has a modern
Go language-like syntax, and a C-like semantics and runtime.

The project is a pure Go project that only  depends on the standard
library. The compiler is written from scratch.

[Try G language in the playground.](https://g.smallrepo.com/play)

[1]: https://github.com/shanhuio/smlvm/wiki/G-Language-Introduction

## Using the `makefile`

The project comes with a `makefile`, which formats the code files,
check lints, check circular dependencies and build tags. Running the
`makefile` requires installing some tools.

```
go get -u shanhu.io/tools/...
go get -u github.com/golang/lint/golint
go get -u github.com/jstemmer/gotags
```

## Comprehensible code

Go language is already fairly readable. To keep the project
modularized and properly structured, we follow two additional rules
on coding:

1. Each file has no more than 300 lines. 80 characters max per line.
2. No circular dependencies among files.

Note that Go language already forbids circular dependencies among
packages.

With no circular dependencies, the project architecture can be
visualized as [a topology map][2]. We find such maps extremely useful
on developing complicated projects.

[2]: https://shanhu.io/smlvm

Similarly, G language projects for the virtual machine also follows
the same rules. In fact, it is enforced by our compiler.
For example, G language's standard library can also be visualized into
into [a topology map][3].

[3]: https://g.smallrepo.com/r/std


## Copyright and License

Copyright by Shanhu Coding Society. Licence Apache.
