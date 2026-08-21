package main

import (
	"fmt"
	"naoslang/parser"
	"naoslang/tokenizer"
	"os"
	"path/filepath"
)

func main() {
	path := "./testFile4.ns"
	name := filepath.Base(path)
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// the lexer is officially finished
	lexer := tokenizer.New(string(data), name)
	parser := parser.New(lexer.Tokenize())
	afile, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Printf("PACKAGE-NAME |\t%v\n", afile.PackageName)
	fmt.Println("------------------------------------------------------------------------------------------------------------------")
	for i, imp := range afile.Imports {
		fmt.Printf("IMPORT %d      \t|\t%v\n", i, imp)
	}
	fmt.Println("------------------------------------------------------------------------------------------------------------------")
	for i, decl := range afile.Decls {
		fmt.Printf("DECLARATION %d \t|\t%v\n", i, decl)
	}
}
