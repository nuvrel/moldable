package astfile

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/nuvrel/moldable/internal/typeast"
)

type ImportSpec struct {
	Path  string
	Alias string
}

type InterfaceSpec struct {
	Name    string
	Methods []*types.Func
}

type File struct {
	packageName string
	imports     []*ImportSpec
	interfaces  []*InterfaceSpec
}

func New(packageName string) *File {
	return &File{
		packageName: packageName,
		imports:     make([]*ImportSpec, 0),
		interfaces:  make([]*InterfaceSpec, 0),
	}
}

func (f *File) HasInterfaces() bool {
	return len(f.interfaces) > 0
}

func (f *File) AddImport(spec *ImportSpec) {
	f.imports = append(f.imports, spec)
}

func (f *File) AddInterface(spec *InterfaceSpec) {
	f.interfaces = append(f.interfaces, spec)
}

func (f *File) Build(qual types.Qualifier) (*ast.File, error) {
	imports := f.buildImports()

	interfaces, err := f.buildInterfaces(qual)
	if err != nil {
		return nil, fmt.Errorf("building interface declarations: %w", err)
	}

	decls := make([]ast.Decl, 0, len(imports)+len(interfaces))

	decls = append(decls, imports...)
	decls = append(decls, interfaces...)

	return &ast.File{
		Name:  ast.NewIdent(f.packageName),
		Decls: decls,
	}, nil
}

func (f *File) buildImports() []ast.Decl {
	if len(f.imports) == 0 {
		return nil
	}

	specs := make([]ast.Spec, 0, len(f.imports))

	for _, is := range f.imports {
		name := ""

		if is.Alias != "" && is.Alias != filepath.Base(is.Path) {
			name = is.Alias
		}

		specs = append(specs, &ast.ImportSpec{
			Name: ast.NewIdent(name),
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("%q", is.Path),
			},
		})
	}

	return []ast.Decl{
		&ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: specs,
		},
	}
}

func (f *File) buildInterfaces(qual types.Qualifier) ([]ast.Decl, error) {
	decls := make([]ast.Decl, 0, len(f.interfaces))

	for _, is := range f.interfaces {
		methods := make([]*ast.Field, 0, len(is.Methods))

		for _, m := range is.Methods {
			expr, err := typeast.Convert(m.Type(), qual)
			if err != nil {
				return nil, fmt.Errorf("converting method %q type: %w", m.Name(), err)
			}

			methods = append(methods, &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent(m.Name()),
				},
				Type: expr,
			})
		}

		decls = append(decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent(is.Name),
					Type: &ast.InterfaceType{
						Methods: &ast.FieldList{
							List: methods,
						},
					},
				},
			},
		})
	}

	return decls, nil
}
