package typeast_test

import (
	"go/ast"
	"go/token"
	"go/types"
	"testing"

	"github.com/nuvrel/moldable/internal/typeast"
	"github.com/stretchr/testify/assert"
)

type dummy struct{}

func (dummy) String() string {
	return ""
}

func (dummy) Underlying() types.Type {
	return nil
}

var (
	_ types.Type = (*dummy)(nil)
)

func TestConvert(t *testing.T) {
	t.Parallel()

	p1 := types.NewPackage("p1", "p1")
	p2 := types.NewPackage("p2", "p2")

	tn1 := types.NewTypeName(token.NoPos, p1, "", nil)
	tn2 := types.NewTypeName(token.NoPos, p2, "", nil)

	in := types.NewInterfaceType(nil, nil)

	n1 := types.NewNamed(tn1, in, nil)
	n2 := types.NewNamed(tn2, in, nil)

	successes := []struct {
		name  string
		input types.Type
		want  ast.Expr
	}{
		{
			name:  "basic",
			input: types.Typ[types.Int],
			want:  ast.NewIdent("int"),
		},
		{
			name:  "array",
			input: types.NewArray(n1, 10),
			want: &ast.ArrayType{
				Len: &ast.BasicLit{
					Kind:  token.INT,
					Value: "10",
				},
				Elt: &ast.SelectorExpr{X: ast.NewIdent("p1"), Sel: ast.NewIdent("")},
			},
		},
		{
			name:  "slice",
			input: types.NewSlice(n1),
			want: &ast.ArrayType{
				Elt: &ast.SelectorExpr{X: ast.NewIdent("p1"), Sel: ast.NewIdent("")},
			},
		},
		{
			name:  "map",
			input: types.NewMap(n1, n2),
			want: &ast.MapType{
				Key:   &ast.SelectorExpr{X: ast.NewIdent("p1"), Sel: ast.NewIdent("")},
				Value: &ast.SelectorExpr{X: ast.NewIdent("p2"), Sel: ast.NewIdent("")},
			},
		},
		{
			name:  "chan bidir",
			input: types.NewChan(types.SendOnly|types.RecvOnly, n1),
			want: &ast.ChanType{
				Dir:   ast.SEND | ast.RECV,
				Value: &ast.SelectorExpr{X: ast.NewIdent("p1"), Sel: ast.NewIdent("")},
			},
		},
		{
			name:  "chan send",
			input: types.NewChan(types.SendOnly, n1),
			want: &ast.ChanType{
				Dir:   ast.SEND,
				Value: &ast.SelectorExpr{X: ast.NewIdent("p1"), Sel: ast.NewIdent("")},
			},
		},
		{
			name:  "chan recv",
			input: types.NewChan(types.RecvOnly, n1),
			want: &ast.ChanType{
				Dir:   ast.RECV,
				Value: &ast.SelectorExpr{X: ast.NewIdent("p1"), Sel: ast.NewIdent("")},
			},
		},
		{
			name:  "pointer",
			input: types.NewPointer(n1),
			want: &ast.StarExpr{
				X: &ast.SelectorExpr{X: ast.NewIdent("p1"), Sel: ast.NewIdent("")},
			},
		},
		{
			name:  "empty interface",
			input: in,
			want:  ast.NewIdent("any"),
		},
		{
			name: "signature with params",
			input: func() *types.Signature {
				return types.NewSignatureType(
					nil,
					nil,
					nil,
					types.NewTuple(types.NewVar(token.NoPos, nil, "", n1)),
					nil,
					false,
				)
			}(),
			want: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("")},
							Type:  &ast.SelectorExpr{X: ast.NewIdent("p1"), Sel: ast.NewIdent("")},
						},
					},
				},
				Results: &ast.FieldList{},
			},
		},
		{
			name: "signature with results",
			input: func() *types.Signature {
				return types.NewSignatureType(
					nil,
					nil,
					nil,
					nil,
					types.NewTuple(types.NewVar(token.NoPos, nil, "", n1)),
					false,
				)
			}(),
			want: &ast.FuncType{
				Params: &ast.FieldList{},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("")},
							Type:  &ast.SelectorExpr{X: ast.NewIdent("p1"), Sel: ast.NewIdent("")},
						},
					},
				},
			},
		},
		{
			name: "variadic signature",
			input: func() *types.Signature {
				return types.NewSignatureType(
					nil,
					nil,
					nil,
					types.NewTuple(types.NewVar(token.NoPos, nil, "", types.NewSlice(n1))),
					nil,
					true,
				)
			}(),
			want: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{ast.NewIdent("")},
							Type: &ast.Ellipsis{
								Elt: &ast.SelectorExpr{X: ast.NewIdent("p1"), Sel: ast.NewIdent("")},
							},
						},
					},
				},
				Results: &ast.FieldList{},
			},
		},
		{
			name:  "named",
			input: n1,
			want:  &ast.SelectorExpr{X: ast.NewIdent("p1"), Sel: ast.NewIdent("")},
		},
		{
			name:  "named without qualifier",
			input: types.NewNamed(types.NewTypeName(token.NoPos, nil, "", nil), in, nil),
			want:  &ast.Ident{Name: ""},
		},
		{
			name:  "alias",
			input: types.NewAlias(tn1, n1),
			want:  &ast.SelectorExpr{X: ast.NewIdent("p1"), Sel: ast.NewIdent("")},
		},
		{
			name:  "alias without qualifier",
			input: types.NewAlias(types.NewTypeName(token.NoPos, nil, "", nil), in),
			want:  &ast.Ident{Name: ""},
		},
		{
			name:  "type param",
			input: types.NewTypeParam(tn1, nil),
			want:  &ast.Ident{Name: ""},
		},
		{
			name:  "union",
			input: types.NewUnion([]*types.Term{types.NewTerm(false, n1), types.NewTerm(false, n2)}),
			want: &ast.BinaryExpr{
				X:  &ast.SelectorExpr{X: ast.NewIdent("p1"), Sel: ast.NewIdent("")},
				Op: token.OR,
				Y:  &ast.SelectorExpr{X: ast.NewIdent("p2"), Sel: ast.NewIdent("")},
			},
		},
		{
			name:  "union with tilde",
			input: types.NewUnion([]*types.Term{types.NewTerm(true, types.Typ[types.Int])}),
			want: &ast.UnaryExpr{
				Op: token.TILDE,
				X:  ast.NewIdent("int"),
			},
		},
	}

	for _, s := range successes {
		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			got, err := typeast.Convert(s.input, func(pkg *types.Package) string {
				return pkg.Name()
			})

			assert.NoError(t, err)
			assert.Equal(t, s.want, got)
		})
	}

	failures := []struct {
		name  string
		input types.Type
		err   string
	}{
		{
			name:  "nil",
			input: nil,
			err:   "cannot convert nil type",
		},
		{
			name:  "unsupported",
			input: &dummy{},
			err:   "unsupported type",
		},
		{
			name:  "struct",
			input: types.NewStruct(nil, nil),
			err:   "struct is not supported",
		},
		{
			name:  "array with unsupported type",
			input: types.NewArray(&dummy{}, 10),
			err:   "converting array elem",
		},
		{
			name:  "slice with unsupported type",
			input: types.NewSlice(&dummy{}),
			err:   "converting slice elem",
		},
		{
			name:  "map with unsupported key",
			input: types.NewMap(&dummy{}, types.Typ[types.Int]),
			err:   "converting map key",
		},
		{
			name:  "map with unsupported value",
			input: types.NewMap(types.Typ[types.Int], &dummy{}),
			err:   "converting map value",
		},
		{
			name:  "chan with unsupported type",
			input: types.NewChan(types.SendOnly, &dummy{}),
			err:   "converting chan elem",
		},
		{
			name:  "pointer with unsupported type",
			input: types.NewPointer(&dummy{}),
			err:   "converting pointer elem",
		},
		{
			name:  "non-empty interface",
			input: types.NewInterfaceType(nil, []types.Type{&dummy{}}),
			err:   "non-empty interface is not supported",
		},
		{
			name: "signature with unsupported type in type params",
			input: func() *types.Signature {
				return types.NewSignatureType(
					nil,
					nil,
					nil,
					types.NewTuple(types.NewVar(token.NoPos, nil, "", &dummy{})),
					nil,
					false,
				)
			}(),
			err: "converting signature params",
		},
		{
			name: "signature with unsupported type in results",
			input: func() *types.Signature {
				return types.NewSignatureType(
					nil,
					nil,
					nil,
					nil,
					types.NewTuple(types.NewVar(token.NoPos, nil, "", &dummy{})),
					false,
				)
			}(),
			err: "converting signature results",
		},
		{
			name: "union with unsupported [0] term type",
			input: func() types.Type {
				return types.NewUnion([]*types.Term{
					types.NewTerm(false, &dummy{}),
				})
			}(),
			err: "converting union term 0",
		},
		{
			name: "union with unsupported [1] term type",
			input: func() types.Type {
				return types.NewUnion([]*types.Term{
					types.NewTerm(false, types.Typ[types.Int]),
					types.NewTerm(false, &dummy{}),
				})
			}(),
			err: "converting union term 1",
		},
	}

	for _, f := range failures {
		t.Run(f.name, func(t *testing.T) {
			t.Parallel()

			_, err := typeast.Convert(f.input, nil)

			assert.ErrorContains(t, err, f.err)
		})
	}
}

func TestConvertTypeParams(t *testing.T) {
	t.Parallel()

	tn1 := types.NewTypeName(token.NoPos, types.NewPackage("p1", "p1"), "", nil)
	tn2 := types.NewTypeName(token.NoPos, types.NewPackage("p2", "p2"), "", nil)

	n1 := types.NewNamed(tn1, types.NewStruct(nil, nil), nil)
	n2 := types.NewNamed(tn2, types.NewInterfaceType(nil, nil), nil)

	successes := []struct {
		name  string
		input *types.TypeParamList
		want  *ast.FieldList
	}{
		{
			name:  "nil",
			input: nil,
			want:  nil,
		},
		{
			name: "type params",
			// type p1.[p1. p2.] struct{}
			input: func() *types.TypeParamList {
				n1.SetTypeParams([]*types.TypeParam{
					types.NewTypeParam(tn1, n2),
				})

				return n1.TypeParams()
			}(),
			want: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("")},
						Type:  &ast.SelectorExpr{X: ast.NewIdent("p2"), Sel: ast.NewIdent("")},
					},
				},
			},
		},
	}

	for _, s := range successes {
		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			got, err := typeast.ConvertTypeParams(s.input, func(pkg *types.Package) string {
				return pkg.Name()
			})

			assert.NoError(t, err)
			assert.Equal(t, s.want, got)
		})
	}

	failures := []struct {
		name  string
		input *types.TypeParamList
		err   string
	}{
		{
			name: "unsupported type in type params",
			input: func() *types.TypeParamList {
				n1.SetTypeParams([]*types.TypeParam{
					types.NewTypeParam(tn1, &dummy{}),
				})

				return n1.TypeParams()
			}(),
			err: "converting type param constraint 0",
		},
	}

	for _, f := range failures {
		t.Run(f.name, func(t *testing.T) {
			t.Parallel()

			_, err := typeast.ConvertTypeParams(f.input, func(pkg *types.Package) string {
				return pkg.Name()
			})

			assert.ErrorContains(t, err, f.err)
		})
	}
}
