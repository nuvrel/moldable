package structcollector

import "go/types"

type StructSpec struct {
	TypeName   *types.TypeName
	TypeParams *types.TypeParamList
	Methods    []*types.Func
}

func (ss StructSpec) HasMethods() bool {
	return len(ss.Methods) > 0
}

type StructCollector struct {
	structs map[*types.Package][]*StructSpec
}

func New() *StructCollector {
	return &StructCollector{
		structs: make(map[*types.Package][]*StructSpec),
	}
}

func (sc *StructCollector) Collect(pkg *types.Package) []*StructSpec {
	if info, ok := sc.structs[pkg]; ok {
		return info
	}

	structs := sc.analyzePackage(pkg)

	sc.structs[pkg] = structs

	return structs
}

func (sc *StructCollector) analyzePackage(pkg *types.Package) []*StructSpec {
	scope := pkg.Scope()
	names := scope.Names()

	structs := make([]*StructSpec, 0, len(names))

	for _, name := range names {
		if spec := sc.analyzeObject(scope.Lookup(name)); spec != nil {
			structs = append(structs, spec)
		}
	}

	return structs
}

func (sc *StructCollector) analyzeObject(obj types.Object) *StructSpec {
	tn, ok := obj.(*types.TypeName)
	if !ok || !tn.Exported() {
		return nil
	}

	_, ok = tn.Type().Underlying().(*types.Struct)
	if !ok {
		return nil
	}

	tp := sc.collectTypeParams(tn)
	methods := sc.collectMethods(tn)

	return &StructSpec{
		TypeName:   tn,
		TypeParams: tp,
		Methods:    methods,
	}
}

func (sc *StructCollector) collectMethods(tn *types.TypeName) []*types.Func {
	ms := types.NewMethodSet(types.NewPointer(tn.Type()))

	methods := make([]*types.Func, 0, ms.Len())

	for i := range ms.Len() {
		if m, ok := ms.At(i).Obj().(*types.Func); ok && m.Exported() {
			methods = append(methods, m)
		}
	}

	return methods
}

func (sc *StructCollector) collectTypeParams(tn *types.TypeName) *types.TypeParamList {
	if named, ok := tn.Type().(*types.Named); ok {
		return named.TypeParams()
	}

	return nil
}
