Hatch, in Go
============

The goal of this work is to produce a Lisp compiler that targets the
Z80/8080 instruction set. It's aim is not for general purpose use but as
a learning aid for compiler construction.

This work is a continuation of my initial work on
[Hatch](https://github.com/bramwelt/hatch) but rewritten in Go, and
extended to be Lisp and do code-generation it is still a simple LL(1)
recursive descent compiler.

Building
--------
The code can be built with `go build` or ran directly with `go run *.go`
