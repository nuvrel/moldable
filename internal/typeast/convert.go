package typeast

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
)

// TODO(calmondev): return typed errors to reuse in tests instead of relying on
// text alone
func Convert(typ types.Type, qual types.Qualifier) (ast.Expr, error) {
	if typ == nil {
		return nil, fmt.Errorf("cannot convert nil type")
	}

	return convert(typ, qual)
}

func convert(typ types.Type, qual types.Qualifier) (ast.Expr, error) {
	switch t := typ.(type) {
	case *types.Basic:
		return convertBasic(t)
	case *types.Struct:
		return convertStruct()
	case *types.Array:
		return convertArray(t, qual)
	case *types.Slice:
		return convertSlice(t, qual)
	case *types.Map:
		return convertMap(t, qual)
	case *types.Chan:
		return convertChan(t, qual)
	case *types.Pointer:
		return convertPointer(t, qual)
	case *types.Interface:
		return convertInterface(t)
	case *types.Signature:
		return convertSignature(t, qual)
	case *types.Named:
		return convertNamed(t, qual)
	case *types.Alias:
		return convertAlias(t, qual)
	case *types.TypeParam:
		return convertTypeParam(t)
	case *types.Union:
		return convertUnion(t, qual)
	default:
		return nil, fmt.Errorf("unsupported type %T", t)
	}
}

func convertBasic(b *types.Basic) (ast.Expr, error) {
	return ast.NewIdent(b.Name()), nil
}

func convertStruct() (ast.Expr, error) {
	return nil, fmt.Errorf("struct is not supported")
}

func convertArray(arr *types.Array, qual types.Qualifier) (ast.Expr, error) {
	el, err := convert(arr.Elem(), qual)
	if err != nil {
		return nil, fmt.Errorf("converting array elem: %w", err)
	}

	return &ast.ArrayType{
		Len: &ast.BasicLit{
			Kind:  token.INT,
			Value: strconv.FormatInt(arr.Len(), 10),
		},
		Elt: el,
	}, nil
}

func convertSlice(s *types.Slice, qual types.Qualifier) (ast.Expr, error) {
	el, err := convert(s.Elem(), qual)
	if err != nil {
		return nil, fmt.Errorf("converting slice elem: %w", err)
	}

	return &ast.ArrayType{
		Elt: el,
	}, nil
}

func convertMap(m *types.Map, qual types.Qualifier) (ast.Expr, error) {
	key, err := convert(m.Key(), qual)
	if err != nil {
		return nil, fmt.Errorf("converting map key: %w", err)
	}

	val, err := convert(m.Elem(), qual)
	if err != nil {
		return nil, fmt.Errorf("converting map value: %w", err)
	}

	return &ast.MapType{
		Key:   key,
		Value: val,
	}, nil
}

func convertChan(ch *types.Chan, qual types.Qualifier) (ast.Expr, error) {
	el, err := convert(ch.Elem(), qual)
	if err != nil {
		return nil, fmt.Errorf("converting chan elem: %w", err)
	}

	var dir ast.ChanDir = ast.SEND | ast.RECV

	switch ch.Dir() {
	case types.SendOnly:
		dir = ast.SEND
	case types.RecvOnly:
		dir = ast.RECV
	}

	return &ast.ChanType{
		Dir:   dir,
		Value: el,
	}, nil
}

func convertPointer(ptr *types.Pointer, qual types.Qualifier) (ast.Expr, error) {
	el, err := convert(ptr.Elem(), qual)
	if err != nil {
		return nil, fmt.Errorf("converting pointer elem: %w", err)
	}

	return &ast.StarExpr{
		X: el,
	}, nil
}

func convertInterface(iface *types.Interface) (ast.Expr, error) {
	if !iface.Empty() {
		return nil, fmt.Errorf("non-empty interface is not supported")
	}

	return ast.NewIdent("any"), nil
}

func convertSignature(sig *types.Signature, qual types.Qualifier) (ast.Expr, error) {
	params, err := convertTuple(sig.Params(), qual)
	if err != nil {
		return nil, fmt.Errorf("converting signature params: %w", err)
	}

	results, err := convertTuple(sig.Results(), qual)
	if err != nil {
		return nil, fmt.Errorf("converting signature results: %w", err)
	}

	if sig.Variadic() {
		last := params.List[len(params.List)-1]

		last.Type = &ast.Ellipsis{
			Elt: last.Type.(*ast.ArrayType).Elt,
		}
	}

	return &ast.FuncType{
		Params:  params,
		Results: results,
	}, nil
}

func convertTuple(tup *types.Tuple, qual types.Qualifier) (*ast.FieldList, error) {
	if tup == nil {
		return &ast.FieldList{}, nil
	}

	list := make([]*ast.Field, tup.Len())

	for i := range tup.Len() {
		v := tup.At(i)

		typ, err := convert(v.Type(), qual)
		if err != nil {
			return nil, fmt.Errorf("converting tuple field %d: %w", i, err)
		}

		list[i] = &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(v.Name())},
			Type:  typ,
		}
	}

	return &ast.FieldList{
		List: list,
	}, nil
}

func convertNamed(n *types.Named, qual types.Qualifier) (ast.Expr, error) {
	return qualifyTypeName(n.Obj(), qual), nil
}

func convertAlias(a *types.Alias, qual types.Qualifier) (ast.Expr, error) {
	return qualifyTypeName(a.Obj(), qual), nil
}

func qualifyTypeName(tn *types.TypeName, qual types.Qualifier) ast.Expr {
	id := ast.NewIdent(tn.Name())

	if pkg := tn.Pkg(); pkg != nil && qual != nil && qual(pkg) != "" {
		return &ast.SelectorExpr{
			X:   ast.NewIdent(qual(pkg)),
			Sel: id,
		}
	}

	return id
}

func convertTypeParam(tp *types.TypeParam) (ast.Expr, error) {
	return ast.NewIdent(tp.Obj().Name()), nil
}

func convertUnion(u *types.Union, qual types.Qualifier) (ast.Expr, error) {
	expr, err := convertUnionTerm(u.Term(0), qual)
	if err != nil {
		return nil, fmt.Errorf("converting union term 0: %w", err)
	}

	for i := range u.Len() {
		if i == 0 {
			continue
		}

		y, err := convertUnionTerm(u.Term(i), qual)
		if err != nil {
			return nil, fmt.Errorf("converting union term %d: %w", i, err)
		}

		expr = &ast.BinaryExpr{
			X:  expr,
			Op: token.OR,
			Y:  y,
		}
	}

	return expr, nil
}

func convertUnionTerm(t *types.Term, qual types.Qualifier) (ast.Expr, error) {
	expr, err := convert(t.Type(), qual)
	if err != nil {
		return nil, fmt.Errorf("converting union term type: %w", err)
	}

	if t.Tilde() {
		expr = &ast.UnaryExpr{
			Op: token.TILDE,
			X:  expr,
		}
	}

	return expr, nil
}

func ConvertTypeParams(tpl *types.TypeParamList, qual types.Qualifier) (*ast.FieldList, error) {
	if tpl == nil {
		return nil, nil
	}

	list := make([]*ast.Field, tpl.Len())

	for i := range tpl.Len() {
		tp := tpl.At(i)

		typ, err := convert(tp.Constraint(), qual)
		if err != nil {
			return nil, fmt.Errorf("converting type param constraint %d: %w", i, err)
		}

		list[i] = &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(tp.Obj().Name())},
			Type:  typ,
		}
	}

	return &ast.FieldList{
		List: list,
	}, nil
}
