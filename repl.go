/*
SPDX-License-Identifier: Apache-2.0
Copyright (c) 2019 Trevor Bramwell
*/
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const PS1 = "> "

// 'Read, Evaluate, Print, Loop'
func Repl() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(PS1)
		text := Read(reader)
		lineReader := bufio.NewReader(strings.NewReader(text))
		output := Eval(lineReader)
		Print(output)
	}
}

// Read: Return a newline delimited string from the reader
func Read(reader *bufio.Reader) string {
	l, err := reader.ReadString('\n')
	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
		fmt.Println()
		os.Exit(0)
	}
	return l
}

func Print(text string) {
	fmt.Println(text)
}

// Evaluate: Tokenize the stream (AST), parse, and return the result
func Eval(r *bufio.Reader) string {
	tokenizer := Tokenize(r)
	result := Parse(tokenizer)
	return strconv.Itoa(result)
}
