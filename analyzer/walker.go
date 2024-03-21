package analyzer

import (
	"unicode"

	sitter "github.com/smacker/go-tree-sitter"
)

// FuncInfo stores function name and comments associated with it.
//
// comments must be positioned above the function declaration
// and should be sequential and non split by whitespace (e.g. newlines without comment token).
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
	var (
        symbols PublicSymbols
        queue []*sitter.Node
    )

	queue = append(queue, n)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		switch curr.Type() {
		case "function_declaration":
			symbols.Functions = append(symbols.Functions, psc.findPublicFunctions(curr)...)
		case "const_declaration":
			symbols.Constants = append(symbols.Constants, psc.findPublicSymbols(curr)...)
		case "var_declaration":
			symbols.Variables = append(symbols.Variables, psc.findPublicSymbols(curr)...)
		}

		for i := 0; i < int(curr.ChildCount()); i++ {
			queue = append(queue, curr.Child(i))
		}
	}

	return symbols
}

// collectComments collects comments associated with the given node.
func (psc *PublicSymbolsCollector) collectComments(n *sitter.Node) (comments []string) {
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
func (psc *PublicSymbolsCollector) findPublicFunctions(n *sitter.Node) (pubfns []FuncInfo) {
    if n.Type() != "function_declaration" {
        return pubfns
    }

    nn := n.ChildByFieldName("name")
    if nn == nil {
        return pubfns
    }

    ident := nn.Content([]byte(psc.code))
    if !unicode.IsUpper(rune(ident[0])) {
        return pubfns
    }

    comments := psc.collectComments(n)
    pubfns = append(pubfns, FuncInfo{
        Ident:    ident,
        Comments: comments,
    })

    return pubfns
}

// findPublicSymbols traverses the given node and finds public symbols.
//
// Like godoc, autodoc also extract the public constants and variables.
// To find public constants and variables, it looks for const_declaration
// and var_declaration nodes and checks if the symbol name is public
// (starts with an uppercase letter).
func (psc *PublicSymbolsCollector) findPublicSymbols(n *sitter.Node) (pubsyms []string) {
    switch n.Type() {
    case "const_declaration", "var_declaration":
        for i := 0; i < int(n.NamedChildCount()); i++ {
            child := n.NamedChild(i)
            pubsyms = appendPublicSymbols(pubsyms, findPublicSymbolsInNode(child, []byte(psc.code)))
        }
    }

    return pubsyms
}

// findPublicSymbols traverses the given node and finds public symbols.
//
// Like godoc, autodoc also extract the public constants and variables.
// To find public constants and variables, it looks for const_declaration
// and var_declaration nodes and checks if the symbol name is public
// (starts with an uppercase letter).
func findPublicSymbolsInNode(n *sitter.Node, code []byte) (pubsyms []string) {
    for i := 0; i < int(n.NamedChildCount()); i++ {
        nameNode := n.NamedChild(i)
        if nameNode == nil {
            continue
        }

        name := nameNode.Content(code)
        if unicode.IsUpper(rune(name[0])) {
            pubsyms = append(pubsyms, name)
        }
    }
    return pubsyms
}

func appendPublicSymbols(dest, src []string) []string {
    tmp := make([]string, 0, len(dest)+len(src))
    tmp = append(tmp, dest...)
    tmp = append(tmp, src...)
    return tmp
}
