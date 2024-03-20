package analyzer

import (
	"unicode"

	sitter "github.com/smacker/go-tree-sitter"
)

// FuncInfo stores function name and comments associated with it.
//
// comments must be positioned above the function declaration
// and should be sequential.
type FuncInfo struct {
	Ident    string   // function name (identifier)
	Comments []string // comments associated with the function
}

// PublicSymbols stores public symbols (functions, constants, and variables).
type PublicSymbols struct {
	Functions []FuncInfo
	Constants []string
	Variables []string
}

// NewPublicSymbolsCollector creates a new instance of PublicSymbolsCollector.
func NewPublicSymbolsCollector(code string) *PublicSymbolsCollector {
	return &PublicSymbolsCollector{
		code: code,
	}
}

// PublicSymbolsCollector collects public symbols from the given AST.
type PublicSymbolsCollector struct {
	code string
}

// Collect traverses the AST and collects public symbols.
func (psc *PublicSymbolsCollector) Collect(n *sitter.Node) PublicSymbols {
	var symbols PublicSymbols
	var queue []*sitter.Node
	queue = append(queue, n)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		switch current.Type() {
		case "function_declaration":
			symbols.Functions = append(symbols.Functions, psc.findPublicFunctions(current)...)
		case "const_declaration":
			symbols.Constants = append(symbols.Constants, psc.findPublicSymbols(current)...)
		case "var_declaration":
			symbols.Variables = append(symbols.Variables, psc.findPublicSymbols(current)...)
		}

		for i := 0; i < int(current.ChildCount()); i++ {
			queue = append(queue, current.Child(i))
		}
	}

	return symbols
}

// collectComments collects comments associated with the given node.
func (psc *PublicSymbolsCollector) collectComments(n *sitter.Node) []string {
	var comments []string
	prev := n.PrevNamedSibling()
	for prev != nil && prev.Type() == "comment" {
		comments = append([]string{prev.Content([]byte(psc.code))}, comments...)
		prev = prev.PrevNamedSibling()
	}
	return comments
}

// findPublicFunctions traverses the given node and finds public functions.
// To find public functions, it looks for function_declaration nodes
// and checks if the function name is public (starts with an uppercase letter).
func (psc *PublicSymbolsCollector) findPublicFunctions(n *sitter.Node) []FuncInfo {
    var publicFunctions []FuncInfo
    if n.Type() != "function_declaration" {
        return publicFunctions
    }

    funcNameNode := n.ChildByFieldName("name")
    if funcNameNode == nil {
        return publicFunctions
    }

    funcName := funcNameNode.Content([]byte(psc.code))
    if !unicode.IsUpper(rune(funcName[0])) {
        return publicFunctions
    }

    comments := psc.collectComments(n)
    publicFunctions = append(publicFunctions, FuncInfo{
        Ident:    funcName,
        Comments: comments,
    })

    return publicFunctions
}

// findPublicSymbols traverses the given node and finds public symbols.
//
// Like godoc, autodoc also extract the public constants and variables.
// To find public constants and variables, it looks for const_declaration
// and var_declaration nodes and checks if the symbol name is public
// (starts with an uppercase letter).
func (psc *PublicSymbolsCollector) findPublicSymbols(n *sitter.Node) []string {
    var publicSymbols []string

    switch n.Type() {
    case "const_declaration", "var_declaration":
        for i := 0; i < int(n.NamedChildCount()); i++ {
            child := n.NamedChild(i)
            publicSymbols = appendPublicSymbols(publicSymbols, findPublicSymbolsInNode(child, []byte(psc.code)))
        }
    }

    return publicSymbols
}

// findPublicSymbols traverses the given node and finds public symbols.
//
// Like godoc, autodoc also extract the public constants and variables.
// To find public constants and variables, it looks for const_declaration
// and var_declaration nodes and checks if the symbol name is public
// (starts with an uppercase letter).
func findPublicSymbolsInNode(n *sitter.Node, code []byte) []string {
    var publicSymbols []string
    for i := 0; i < int(n.NamedChildCount()); i++ {
        nameNode := n.NamedChild(i)
        if nameNode == nil {
            continue
        }

        name := nameNode.Content(code)
        if unicode.IsUpper(rune(name[0])) {
            publicSymbols = append(publicSymbols, name)
        }
    }
    return publicSymbols
}

func appendPublicSymbols(dest, src []string) []string {
    tmp := make([]string, 0, len(dest)+len(src))
    tmp = append(tmp, dest...)
    tmp = append(tmp, src...)
    return tmp
}
