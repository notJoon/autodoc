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

	code := `
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

		// Public constants
		const (
			PublicConst1 = 1
			publicConst2 = 2
			PublicConst3 = 3
		)

		// Private constants
		const (
			privateConst1 = 4
			privateConst2 = 5
		)

		// Public variables
		var (
			PublicVar1 int
			PublicVar2 string
		)

		// Private variables
		var (
			privateVar1 bool
			privateVar2 float64
		)
	`

	tree, err := parser.ParseCtx(context.Background(), nil, []byte(code))
	if err != nil {
		t.Fatalf("Failed to parse code: %s", err)
	}

	n := tree.RootNode()
	res := WalkTreeNode(n, code)

	expectedFunctions := []string{"Foo"}
	if !reflect.DeepEqual(res.Functions, expectedFunctions) {
		t.Errorf("Expected functions %v, got %v", expectedFunctions, res.Functions)
	}

	expectedConstants := []string{"PublicConst1", "PublicConst3"}
	if !reflect.DeepEqual(res.Constants, expectedConstants) {
		t.Errorf("Expected constants %v, got %v", expectedConstants, res.Constants)
	}

	expectedVariables := []string{"PublicVar1", "PublicVar2"}
	if !reflect.DeepEqual(res.Variables, expectedVariables) {
		t.Errorf("Expected variables %v, got %v", expectedVariables, res.Variables)
	}
}