package typeast

import "go/types"

type Func func(*types.Package)

func Traverse(typ types.Type, fn Func) {
	if typ == nil {
		return
	}

	switch t := typ.(type) {
	case *types.Array:
		Traverse(t.Elem(), fn)
	case *types.Slice:
		Traverse(t.Elem(), fn)
	case *types.Map:
		Traverse(t.Key(), fn)
		Traverse(t.Elem(), fn)
	case *types.Chan:
		Traverse(t.Elem(), fn)
	case *types.Pointer:
		Traverse(t.Elem(), fn)
	case *types.Interface:
		for i := range t.NumEmbeddeds() {
			Traverse(t.EmbeddedType(i), fn)
		}
	case *types.Signature:
		if params := t.Params(); params != nil {
			for i := range params.Len() {
				Traverse(params.At(i).Type(), fn)
			}
		}

		if results := t.Results(); results != nil {
			for i := range results.Len() {
				Traverse(results.At(i).Type(), fn)
			}
		}
	case *types.Named:
		if pkg := t.Obj().Pkg(); pkg != nil {
			fn(pkg)
		}
	case *types.Alias:
		if pkg := t.Obj().Pkg(); pkg != nil {
			fn(pkg)
		}
	case *types.TypeParam:
		Traverse(t.Constraint(), fn)
	}
}

func TraverseTypeParams(tpl *types.TypeParamList, fn Func) {
	if tpl == nil {
		return
	}

	for i := range tpl.Len() {
		Traverse(tpl.At(i), fn)
	}
}

func TraverseFuncs(fns []*types.Func, fn Func) {
	for i := range fns {
		Traverse(fns[i].Type(), fn)
	}
}
