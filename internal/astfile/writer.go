package astfile

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"

	"golang.org/x/tools/imports"
)

type Writer struct{}

func NewWriter() *Writer {
	return &Writer{}
}

func (Writer) Write(fset *token.FileSet, file *ast.File, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	var buf bytes.Buffer

	if err := format.Node(&buf, fset, file); err != nil {
		return fmt.Errorf("formatting node: %w", err)
	}

	formatted, err := imports.Process(path, buf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("processing imports: %w", err)
	}

	if err := os.WriteFile(path, formatted, 0644); err != nil {
		return fmt.Errorf("writing file to disk: %w", err)
	}

	return nil
}
