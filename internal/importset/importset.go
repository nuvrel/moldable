package importset

import (
	"fmt"
	"go/types"
)

type ImportSet struct {
	imports map[string]string
	aliases map[string]bool
}

func New() *ImportSet {
	return &ImportSet{
		imports: make(map[string]string),
		aliases: make(map[string]bool),
	}
}

func (is *ImportSet) Import(pkg *types.Package) {
	path := pkg.Path()

	if _, ok := is.imports[path]; ok {
		return
	}

	alias := func() string {
		base := pkg.Name()

		if !is.aliases[base] {
			return base
		}

		counter := 1

		for {
			alias := fmt.Sprintf("%s%d", base, counter)

			if !is.aliases[alias] {
				return alias
			}

			counter++
		}
	}()

	is.imports[path] = alias
	is.aliases[alias] = true
}

func (is *ImportSet) Qualifier(pkg *types.Package) string {
	if alias, ok := is.imports[pkg.Path()]; ok {
		return alias
	}

	return ""
}

func (is *ImportSet) Imports() map[string]string {
	return is.imports
}
