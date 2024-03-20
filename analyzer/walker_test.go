package analyzer

import (
	"context"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/stretchr/testify/assert"
)

func TestPublicSymbolsCollector(t *testing.T) {
	sourceCode := `
package example

// Add adds two integers.
// 
// It returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}

const Pi = 3.14

var ExportedVar = "exported"
`

	parser := sitter.NewParser()
	defer parser.Close()

	parser.SetLanguage(golang.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, []byte(sourceCode))
	if err != nil {
		t.Fatal(err)
	}
	defer tree.Close()

	collector := NewPublicSymbolsCollector(sourceCode)
	symbols := collector.Collect(tree.RootNode())

	expectedFuncs := []FuncInfo{{Ident: "Add", Comments: []string{
		"// Add adds two integers.",
		"// ",
		"// It returns the sum of two integers.",
	}}}
	expectedConsts := []string{"Pi"}
	expectedVars := []string{"ExportedVar"}

	assert.Equal(t, expectedFuncs, symbols.Functions, "should correctly extract public functions")
	assert.Equal(t, expectedConsts, symbols.Constants, "should correctly extract public constants")
	assert.Equal(t, expectedVars, symbols.Variables, "should correctly extract public variables")
}
