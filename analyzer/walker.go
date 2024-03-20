package analyzer

import (
	"unicode"

	sitter "github.com/smacker/go-tree-sitter"
)

// WalkTreeNode recursively walks through the syntax tree of the source code.
func WalkTreeNode(n *sitter.Node, depth int, code []byte) []string {
	var names []string

	if n.Type() == "function_declaration" {
		name := n.ChildByFieldName("name")
		if name != nil {
			ident := name.Content(code)

			if unicode.IsUpper(rune(ident[0])) {
				names = append(names, ident)
			}
		}
	}

	for i := 0; i < int(n.ChildCount()); i++ {
		child := WalkTreeNode(n.Child(i), depth+1, code)
		names = append(names, child...)
	}

	return names
}
