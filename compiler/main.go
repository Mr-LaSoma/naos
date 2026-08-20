package main

import (
	"naoslang/tokenizer"
	"naoslang/tokenizer/tokens"
	"os"
	"path/filepath"
)

func main() {
	path := "./testFile2.ns"
	name := filepath.Base(path)
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	lexer := tokenizer.New(string(data), name)
	tokens.PrintTokens(lexer.Tokenize(), true)
}
