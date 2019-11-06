/*
SPDX-License-Identifier: Apache-2.0
Copyright (c) 2019 Trevor Bramwell
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

type Args struct {
	Debug    bool
	Tokens   bool
	REPL     bool
	Compile  bool
	Assemble bool
}

var Flags = new(Args)

func init() {
	flag.BoolVar(&Flags.Debug, "debug", false, "enable symbol debugging")
	flag.BoolVar(&Flags.Tokens, "tokens", false, "enable token debugging")
	flag.BoolVar(&Flags.REPL, "repl", false, "start a REPL instead of compiling")
	flag.BoolVar(&Flags.Compile, "compile", false, "compile a file to Z80/8080 assembly")
	flag.BoolVar(&Flags.Assemble, "assemble", false, "compile and assemble a file to Z80/8080 hex code")
}

func Compile(filename string) {
	fp, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	reader := bufio.NewReader(fp)
	var env = &Env{
		symbols: make(map[string]int),
	}
	// TODO: Continue Eval til EOF
	output := Eval(env, reader)
	Print(output)
}

func main() {
	flag.Parse()

	if Flags.REPL {
		Repl()
	}

	if len(flag.Args()) == 0 {
		if Flags.Compile && flag.NArg() < 1 {
			fmt.Fprintf(os.Stderr, "%s\n", "Missing argument to compile")
		}
		flag.Usage()
		os.Exit(2)
	}
	if Flags.Compile {
		f := flag.Arg(0)
		Compile(f)
	}
}
