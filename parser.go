/*
SPDX-License-Identifier: Apache-2.0
Copyright (c) 2019 Trevor Bramwell
*/
package main

type ParseTree struct {
	token       *Token
	left, right *ParseTree
}

type Env struct {
	symbols map[string]int
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
	return nil
}

// Parser needs to construct the grammar it expects
func Parse(env *Env, tz *Tokenizer) int {
	pt := Expression(env, tz)
	/*
		// Debug environment
		for k, v := range env.symbols {
			fmt.Printf("%s: %d\n", k, v)
		}
	*/
	return value(env, pt)
}

// Walk the tree to determine the value
func value(env *Env, pt *ParseTree) int {
	if pt != nil {
		t := pt.token
		if t.Id == OP {
			x := value(env, pt.left)
			y := value(env, pt.right)
			return result(t.Word, x, y)
		} else if t.Id == SYM {
			return env.symbols[t.Word]
		} else if t.Id == NUM {
			return t.Val
		}
	}
	return 0
}

// Given a operation and two numbers, return the result represented by
// the operation
func result(op string, x int, y int) int {
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
