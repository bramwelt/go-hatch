/*
SPDX-License-Identifier: Apache-2.0
Copyright (c) 2019 Trevor Bramwell
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Token struct {
	Id   byte
	Val  int
	Word string
}

type Tokenizer struct {
	tokens   []*Token
	position int
}

const (
	NUM = iota
	OP
	SYM
	DEFINE
	LPEREN
	RPEREN
	EOF
	ILLEG
)

var EOFToken = &Token{EOF, 0, ""}

var TokenId = map[byte]string{
	OP:     "OP",
	NUM:    "NUM",
	LPEREN: "(",
	RPEREN: ")",
	SYM:    "SYM",
	DEFINE: "DEF",
	EOF:    "EOF",
	ILLEG:  "ILLEG",
}

// Output a Token as a String
func (t *Token) String() string {
	switch t.Id {
	case NUM:
		return fmt.Sprintf("%s: %d", TokenId[t.Id], t.Val)
	case SYM:
		return fmt.Sprintf("%s: %s", TokenId[t.Id], t.Word)
	default:
		return TokenId[t.Id]
	}
	return ""
}

// Match the next expected token in the list
//  Returns an error if we run out of tokens or mismatch the next token
func (tz *Tokenizer) Match(symbol byte) bool {
	t := tz.Peek()
	if t.Id == symbol {
		tz.Advance()
		return true
	}
	return false
}

// State if the tokenizer is beyond the list
func (tz *Tokenizer) OutOfBounds() bool {
	return tz.position >= len(tz.tokens)
}

// Advance to the next token in the list
func (tz *Tokenizer) Advance() {
	tz.position++
}

// Show the next token Id, but don't advance position
func (tz *Tokenizer) Peek() *Token {
	if tz.OutOfBounds() {
		return EOFToken
	}
	return tz.tokens[tz.position]
}

// Tokenize returns a Tokenizer containing the list of tokens
func Tokenize(reader *bufio.Reader) *Tokenizer {
	tokenizer := new(Tokenizer)
	tokens := []*Token{}
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				tokens = append(tokens, EOFToken)
				break
			}
		}
		t := GetToken(b, reader)
		if t != nil {
			tokens = append(tokens, t)
		}
	}
	tokenizer.tokens = tokens
	return tokenizer
}

// Print out all the tokens
func (tz *Tokenizer) Print() string {
	var out strings.Builder
	for _, t := range tz.tokens {
		out.WriteString(fmt.Sprintf("%+v\n", t))
	}
	return out.String()
}

// Get a token from the stream - sequential numbers are converted to
// integers, and sequential letters to words
func GetToken(b byte, reader *bufio.Reader) *Token {
	switch b {
	case '+', '-', '*', '/':
		return &Token{OP, 0, string(b)}
	case '(':
		return &Token{LPEREN, 0, ""}
	case ')':
		return &Token{RPEREN, 0, ""}
	case ' ', '\t', '\n':
		return nil
	}
	switch {
	case isNumber(b):
		i := getNumber(b, reader)
		return &Token{NUM, i, ""}
	case isLetter(b):
		b := getWord(b, reader)
		switch b {
		case "define":
			return &Token{DEFINE, 0, b}
		}
		return &Token{SYM, 0, b}
	}
	return &Token{ILLEG, 0, ""}
}

// Convert a stream of numbers to an integer
func getNumber(b byte, reader *bufio.Reader) int {
	i := int(b - '0')
	for {
		a, _ := reader.Peek(1)
		if isNumber(a[0]) {
			i = i*10 + int(a[0]-'0')
			reader.ReadByte()
		} else {
			break
		}
	}
	return i
}

// Convert a stream of letters to a string
func getWord(b byte, reader *bufio.Reader) string {
	var word = []byte{b}
	for {
		a, _ := reader.Peek(1)
		if isLetter(a[0]) {
			word = append(word, a[0])
			reader.ReadByte()
		} else {
			break
		}
	}
	return string(word)
}

// Checks if a byte is in the ASCII letter range
func isLetter(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z')
}

// Checks if a byte is in the ASCII number range
func isNumber(b byte) bool {
	return '0' <= b && b <= '9'
}
