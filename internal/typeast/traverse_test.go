package typeast_test

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/nuvrel/moldable/internal/typeast"
	"github.com/stretchr/testify/assert"
)

func TestTraverse(t *testing.T) {
	t.Parallel()

	tn1 := types.NewTypeName(token.NoPos, types.NewPackage("p1", "p1"), "", nil)
	tn2 := types.NewTypeName(token.NoPos, types.NewPackage("p2", "p2"), "", nil)

	in := types.NewInterfaceType(nil, nil)

	n1 := types.NewNamed(tn1, in, nil)
	n2 := types.NewNamed(tn2, in, nil)

	cases := []struct {
		name  string
		input types.Type
		want  []string
	}{
		{
			name:  "nil",
			input: nil,
			want:  nil,
		},
		{
			name: "array",
			// [10]p1.
			input: types.NewArray(n1, 10),
			want:  []string{"p1"},
		},
		{
			name: "slice",
			// []p1.
			input: types.NewSlice(n1),
			want:  []string{"p1"},
		},
		{
			name: "map",
			// map[p1.]p2.
			input: types.NewMap(n1, n2),
			want:  []string{"p1", "p2"},
		},
		{
			name: "chan",
			// chan<- p1.
			input: types.NewChan(types.SendOnly, n1),
			want:  []string{"p1"},
		},
		{
			name: "pointer",
			// *p1.
			input: types.NewPointer(n1),
			want:  []string{"p1"},
		},
		{
			name: "interface",
			// interface { p1. }
			input: types.NewInterfaceType(nil, []types.Type{n1}),
			want:  []string{"p1"},
		},
		{
			name: "signature",
			// func (p1.) (p2.)
			input: func() *types.Signature {
				v1 := types.NewVar(token.NoPos, nil, "", n1)
				v2 := types.NewVar(token.NoPos, nil, "", n2)

				sig := types.NewSignatureType(
					nil,
					nil,
					nil,
					types.NewTuple(v1),
					types.NewTuple(v2),
					false,
				)

				return sig
			}(),
			want: []string{"p1", "p2"},
		},
		{
			name: "named",
			// p1.
			input: n1,
			want:  []string{"p1"},
		},
		{
			name: "alias",
			// p1.
			input: types.NewAlias(tn1, n1),
			want:  []string{"p1"},
		},
		{
			name: "type param",
			// [p1. p2.]
			input: types.NewTypeParam(tn1, n2),
			want:  []string{"p2"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := []string(nil)

			typeast.Traverse(c.input, func(pkg *types.Package) {
				got = append(got, pkg.Path())
			})

			assert.Equal(t, c.want, got)
		})
	}
}

func TestTraverseTypeParams(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input *types.TypeParamList
		want  []string
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
				tn1 := types.NewTypeName(token.NoPos, types.NewPackage("p1", "p1"), "", nil)
				tn2 := types.NewTypeName(token.NoPos, types.NewPackage("p2", "p2"), "", nil)

				n1 := types.NewNamed(tn1, types.NewStruct(nil, nil), nil)
				n2 := types.NewNamed(tn2, types.NewInterfaceType(nil, nil), nil)

				n1.SetTypeParams([]*types.TypeParam{
					types.NewTypeParam(tn1, n2),
				})

				return n1.TypeParams()
			}(),
			want: []string{"p2"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			got := []string(nil)

			typeast.TraverseTypeParams(c.input, func(pkg *types.Package) {
				got = append(got, pkg.Path())
			})

			assert.Equal(t, c.want, got)
		})
	}
}

func TestTraverseFuncs(t *testing.T) {
	t.Parallel()

	tn1 := types.NewTypeName(token.NoPos, types.NewPackage("p1", "p1"), "", nil)
	tn2 := types.NewTypeName(token.NoPos, types.NewPackage("p2", "p2"), "", nil)

	in := types.NewInterfaceType(nil, nil)

	n1 := types.NewNamed(tn1, in, nil)
	n2 := types.NewNamed(tn2, in, nil)

	// func (p1.) p2.
	m1 := types.NewFunc(
		token.NoPos,
		nil,
		"",
		types.NewSignatureType(
			nil,
			nil,
			nil,
			types.NewTuple(types.NewVar(token.NoPos, nil, "", n1)),
			types.NewTuple(types.NewVar(token.NoPos, nil, "", n2)),
			false,
		),
	)

	got := []string(nil)

	typeast.TraverseFuncs([]*types.Func{m1}, func(pkg *types.Package) {
		got = append(got, pkg.Path())
	})

	assert.Equal(t, []string{"p1", "p2"}, got)
}
