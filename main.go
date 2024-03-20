package main

import (
	"autodoc/analyzer"
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

func main() {
	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())

	code := []byte(`
	// top level coment
	const (
		String = "string"
	)

	// This is a comment
	// This is another comment
	func main() {
		fmt.Println("Hello, World!")
	}

	func Foo() {
		fmt.Println("public")
	}

	func Foo2() {
		fmt.Println("public")
	}

	func bar() {
		fmt.Println("private")
	}
	`)

	tree, err := parser.ParseCtx(context.Background(), nil, code)
	if err != nil {
		panic(err)
	}

	n := tree.RootNode()
	res := analyzer.WalkTreeNode(n, string(code))
	for _, r := range res.Functions {
		println(r)
	}
}
