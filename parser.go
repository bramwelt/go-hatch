/*
SPDX-License-Identifier: Apache-2.0
Copyright (c) 2019 Trevor Bramwell
*/
package main

type ParseTree struct {
	token       *Token
	left, right *ParseTree
}

// Given a list of tokens return an AST
func Expression(tz *Tokenizer) *ParseTree {
	t := tz.Peek()
	if tz.Match(NUM) {
		return &ParseTree{token: t, left: nil, right: nil}
	} else if tz.Match(LPEREN) {
		t = tz.Peek()
		if tz.Match(OP) {
			left := Expression(tz)
			right := Expression(tz)
			tz.Match(RPEREN)
			return &ParseTree{token: t, left: left, right: right}
		}
	}
	return nil
}

// Parser needs to construct the grammar it expects
func Parse(tz *Tokenizer) int {
	pt := Expression(tz)
	return value(pt)
}

// Walk the tree to determine the value
func value(pt *ParseTree) int {
	if pt != nil {
		t := pt.token
		if t.Id == OP {
			x := value(pt.left)
			y := value(pt.right)
			return result(t.Word, x, y)
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
