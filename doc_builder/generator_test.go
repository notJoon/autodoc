package docbuilder

import (
	_ "os"
	_ "path/filepath"
	"testing"

	"autodoc/analyzer"
)

func TestConvertCommentsToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		funcInfo analyzer.FuncInfo
		expected string
	}{
		{
			name: "Single line comment",
			funcInfo: analyzer.FuncInfo{
				Ident:    "Add",
				Comments: []string{"// Add adds two integers."},
			},
			expected: "Add adds two integers.\n\n```go\nAdd\n```\n",
		},
		{
			name: "Multiline comment",
			funcInfo: analyzer.FuncInfo{
				Ident: "Add",
				Comments: []string{
					"// Add adds two integers.",
					"// ",
					"// It returns the sum of two integers.",
				},
			},
			expected: "Add adds two integers.\n\nIt returns the sum of two integers.\n\n```go\nAdd\n```\n",
		},
		{
			name: "Comment with leading space",
			funcInfo: analyzer.FuncInfo{
				Ident:    "Add",
				Comments: []string{"//   Add adds two integers."},
			},
			expected: "Add adds two integers.\n\n```go\nAdd\n```\n",
		},
		{
			name: "No comments",
			funcInfo: analyzer.FuncInfo{
				Ident:    "Add",
				Comments: nil,
			},
			expected: "\n\n```go\nAdd\n```\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToMarkDown(tt.funcInfo)
			if result != tt.expected {
				t.Errorf("(%s) Expected: %q, Got: %q", tt.name, tt.expected, result)
			}
		})
	}
}

// func TestWrite(t *testing.T) {
// 	tempDir, err := os.MkdirTemp("", "test")
// 	if err != nil {
// 		t.Fatalf("Failed to create temporary directory: %v", err)
// 	}
// 	defer os.RemoveAll(tempDir)

// 	tests := []struct {
// 		name        string
// 		packageName string
// 		mdContent   string
// 		wantErr     bool
// 	}{
// 		{
// 			name:        "Write valid markdown",
// 			packageName: "mypackage",
// 			mdContent:   "# My Package\n\nThis is a package description.\n\n```\nMyFunc\n```\n",
// 			wantErr:     false,
// 		},
// 		{
// 			name:        "Write empty markdown",
// 			packageName: "emptypackage",
// 			mdContent:   "",
// 			wantErr:     false,
// 		},
// 		{
// 			name:        "Write markdown with invalid package name",
// 			packageName: "invalid/package/name",
// 			mdContent:   "# Invalid Package\n\nThis is an invalid package.\n",
// 			wantErr:     true,
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			err := Write(filepath.Join(tempDir, tc.packageName), tc.mdContent)

// 			if tc.wantErr && err == nil {
// 				t.Errorf("Expected an error, but got none")
// 			} else if !tc.wantErr && err != nil {
// 				t.Errorf("Unexpected error: %v", err)
// 			} else if !tc.wantErr {
// 				mdPath := filepath.Join(tempDir, tc.packageName, "doc.md")
// 				content, err := os.ReadFile(mdPath)
// 				if err != nil {
// 					t.Errorf("Failed to read markdown file: %v", err)
// 				}
// 				if string(content) != tc.mdContent {
// 					t.Errorf("Unexpected markdown content.\nExpected: %q\nGot: %q", tc.mdContent, string(content))
// 				}
// 			}
// 		})
// 	}
// }