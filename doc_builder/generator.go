package docbuilder

import (
	"autodoc/analyzer"
	"os"
	"path/filepath"
	"strings"
)

func ToMarkDown(funcInfo analyzer.FuncInfo) string {
	var sb strings.Builder

	cleanedComments := make([]string, 0, len(funcInfo.Comments))
	for _, comment := range funcInfo.Comments {
		cleanedComment := strings.TrimPrefix(comment, "//")
		cleanedComment = strings.TrimSpace(cleanedComment)
		cleanedComments = append(cleanedComments, cleanedComment)
	}

	doc := strings.Join(cleanedComments, "\n")

	// convert comments to markdown
	sb.WriteString(doc)
	sb.WriteString("\n\n")
	sb.WriteString("```go\n") // start code block for function signature
	sb.WriteString(funcInfo.Ident)
	sb.WriteString("\n```\n")

	return sb.String()
}

func Write(pkg, content string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, pkg, ".md")

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// write content to file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
