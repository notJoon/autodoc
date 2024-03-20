package analyzer

import (
	"unicode"

	sitter "github.com/smacker/go-tree-sitter"
)

type FunctionInfo struct {
	Ident 	 string 	// function name (identifier)
	Comments []string 	// comments associated with the function
}

type PublicSymbols struct {
    Functions []string
    Constants []string
    Variables []string
}

// WalkTreeNode recursively traverses the given node and collects the names of public functions, constants, and variables.
func WalkTreeNode(n *sitter.Node, code string) PublicSymbols {
    var symbols PublicSymbols
    var queue []*sitter.Node
    queue = append(queue, n)

    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]

        if current.Type() == "function_declaration" {
            symbols.Functions = appendPublicSymbols(symbols.Functions, findPublicFunctions(current, code))
        } else if current.Type() == "const_declaration" {
            symbols.Constants = appendPublicSymbols(symbols.Constants, findPublicConstants(current, code))
        } else if current.Type() == "var_declaration" {
            symbols.Variables = appendPublicSymbols(symbols.Variables, findPublicVariables(current, code))
        }

        for i := 0; i < int(current.ChildCount()); i++ {
            queue = append(queue, current.Child(i))
        }
    }

    return symbols
}

func findPublicFunctions(n *sitter.Node, code string) []string {
    var publicFunctions []string
    funcNameNode := n.ChildByFieldName("name")
    if funcNameNode != nil {
        funcName := funcNameNode.Content([]byte(code))
        if unicode.IsUpper(rune(funcName[0])) {
            publicFunctions = append(publicFunctions, funcName)
        }
    }
    return publicFunctions
}

func findPublicConstants(n *sitter.Node, code string) []string {
    var publicConstants []string
    for i := 0; i < int(n.NamedChildCount()); i++ {
        constSpec := n.NamedChild(i)
        for j := 0; j < int(constSpec.NamedChildCount()); j++ {
            constNameNode := constSpec.NamedChild(j)
            if constNameNode != nil {
                constName := constNameNode.Content([]byte(code))
                if unicode.IsUpper(rune(constName[0])) {
                    publicConstants = append(publicConstants, constName)
                }
            }
        }
    }
    return publicConstants
}

func findPublicVariables(n *sitter.Node, code string) []string {
    var publicVariables []string
    for i := 0; i < int(n.NamedChildCount()); i++ {
        varSpec := n.NamedChild(i)
        for j := 0; j < int(varSpec.NamedChildCount()); j++ {
            varNameNode := varSpec.NamedChild(j)
            if varNameNode != nil {
                varName := varNameNode.Content([]byte(code))
                if unicode.IsUpper(rune(varName[0])) {
                    publicVariables = append(publicVariables, varName)
                }
            }
        }
    }
    return publicVariables
}

func appendPublicSymbols(dest, src []string) []string {
    tmp := make([]string, 0, len(dest)+len(src))
    tmp = append(tmp, dest...)
    tmp = append(tmp, src...)
    return tmp
}
