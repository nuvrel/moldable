package pkgload

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Loader struct {
	packages map[string]*types.Package
}

func NewLoader() *Loader {
	return &Loader{
		packages: make(map[string]*types.Package),
	}
}

func (l *Loader) Load(paths []string) error {
	if len(paths) == 0 {
		return nil
	}

	pkgs, err := packages.Load(
		&packages.Config{
			Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo,
		},
		paths...,
	)
	if err != nil {
		return fmt.Errorf("loading go packages: %w", err)
	}

	if len(pkgs) == 0 {
		return fmt.Errorf("no packages found for paths: %v", paths)
	}

	for _, p := range pkgs {
		if len(p.Errors) > 0 {
			return fmt.Errorf("errors loading package %q: %v", p.PkgPath, p.Errors)
		}

		l.packages[p.PkgPath] = p.Types
	}

	for _, path := range paths {
		if _, ok := l.packages[path]; !ok {
			return fmt.Errorf("requested package %q was not loaded", path)
		}
	}

	return nil
}

func (l *Loader) Package(path string) (*types.Package, error) {
	pkg, ok := l.packages[path]
	if !ok {
		return nil, fmt.Errorf("package %q not found", path)
	}

	return pkg, nil
}
