/*
SPDX-License-Identifier: Apache-2.0
Copyright (c) 2019 Trevor Bramwell
*/
package main

import (
	"fmt"
	"strconv"
	"strings"
)

const asmHeader = `.nolist
#include "ti83plus.inc"
.list
.org 9D93h
.db $BB,$6D
	bcall(_ClrLCDFull)
	ld a,0
	ld (CURCOL),a
	ld (CURROW),a
`

const asmFooter = `	bcall(_DispHL)
	bcall(_NewLine)
	ret
.end
.end
`

type ParseTree struct {
	token       *Token
	left, right *ParseTree
}

var Register = map[int]string{
	0: "hl",
	1: "bc",
}

type Env struct {
	symbols  map[string]int
	numsyms  int
	register int
}

// Given a list of tokens return an AST
func Expression(env *Env, tz *Tokenizer) *ParseTree {
	t := tz.Peek()
	if tz.Match(NUM) || tz.Match(SYM) {
		return &ParseTree{token: t, left: nil, right: nil}
	} else if tz.Match(LPEREN) {
		if t = tz.Peek(); tz.Match(OP) {
			left := Expression(env, tz)
			right := Expression(env, tz)
			tz.Match(RPEREN)
			return &ParseTree{token: t, left: left, right: right}
		} else if tz.Match(DEF) {
			if sym := tz.Peek(); tz.Match(SYM) {
				if num := tz.Peek(); tz.Match(NUM) {
					env.symbols[sym.Word] = num.Val
					return &ParseTree{token: sym, left: nil, right: nil}
				}
			}
		}
	}
	// TODO: Handle RPEREN and returning Errors as Expression only
	// constructs ParseTrees for single lines currently
	return nil
}

// Parser needs to construct the grammar it expects
func Parse(env *Env, tz *Tokenizer) string {
	pt := Expression(env, tz)
	if Flags.Debug {
		// Debug symbols and tree
		for k, v := range env.symbols {
			fmt.Printf("sym %s: %d\n", k, v)
		}
		fmt.Printf("%+v", pt)
	}
	if Flags.REPL {
		return strconv.Itoa(IRValue(env, pt))
	}
	if Flags.Compile {
		stb := new(strings.Builder)
		stb.WriteString(asmHeader)
		ASM(env, pt, stb)
		stb.WriteString(asmFooter)
		return stb.String()
	}
	return ""
}

// IRValue Walk the tree to determine the value
func IRValue(env *Env, pt *ParseTree) int {
	if pt != nil {
		t := pt.token
		if t.Id == OP {
			x := IRValue(env, pt.left)
			y := IRValue(env, pt.right)
			return IRResult(t.Word, x, y)
		} else if t.Id == SYM {
			return env.symbols[t.Word]
		} else if t.Id == NUM {
			return t.Val
		}
	}
	return 0
}

// ASM Walk the tree and ouput a string of assembly
func ASM(env *Env, pt *ParseTree, stb *strings.Builder) *Token {
	if pt == nil {
		return nil
	}
	t := pt.token
	switch t.Id {
	case OP:
		x := ASM(env, pt.left, stb)
		y := ASM(env, pt.right, stb)
		switch t.Word {
		case "+":
			stb.WriteString(fmt.Sprintf("\tld hl,%d\n\tld bc,%d\n", x.Val, y.Val))
			stb.WriteString(fmt.Sprintf("\tadd hl,bc\n"))
		case "-":
			stb.WriteString(fmt.Sprintf("\tsub %d,%d\n", y.Val, x.Val))
		case "*":
			stb.WriteString("\tmul")
		case "/":
			stb.WriteString("\tdiv")
		}
		return t
	case SYM:
		// Increment symbol pointer?
		//sym := strconv.Itoa(env.symbols[t.Word])
		stb.WriteString(fmt.Sprintf("\tld %v,%d\n", Register[env.register], env.symbols[t.Word]))
		env.register++
		return t
	case NUM:
		// TODO: Store in data?
		/*
			stb.WriteString(fmt.Sprintf("ld %d,hl%d\n", t.Val, env.register))
			env.register++
		*/
		return t
	case LPEREN:
		return ASM(env, pt.left, stb)
	default:
		return t
	}
	return nil
}

// Given a operation and two numbers, return the result represented by
// the operation
func IRResult(op string, x int, y int) int {
	switch op {
	case "+":
		return x + y
	case "-":
		return x - y
	case "*":
		return x * y
	case "/":
		return x / y
	}
	return 0
}
