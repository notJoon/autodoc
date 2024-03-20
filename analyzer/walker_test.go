package analyzer

import (
	"context"
	"reflect"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

func TestWalkTree(t *testing.T) {
	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())

	code := []byte(`
		// Some comment
		func main() {
			fmt.Println("Hello, World!")
		}

		// Public function comment
		func Foo() {
			fmt.Println("public")
		}

		// Private function comment
		func bar() {
			fmt.Println("private")
		}
	`)

	tree, err := parser.ParseCtx(context.Background(), nil, code)
	if err != nil {
		t.Fatalf("Failed to parse code: %s", err)
	}

	n := tree.RootNode()
	publicFunctions := WalkTreeNode(n, 0, code)

	expected := []string{"Foo"}
	if !reflect.DeepEqual(publicFunctions, expected) {
		t.Fatalf("Expected %v, got %v", expected, publicFunctions)
	}
}