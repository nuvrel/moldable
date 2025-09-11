package app

import (
	"errors"
	"fmt"
	"go/token"
	"strings"
	"unicode"
)

type Config struct {
	Output   Output    `koanf:"output"`
	Packages []Package `koanf:"packages"`
}

func (c Config) Check() error {
	if err := c.Output.check(); err != nil {
		return fmt.Errorf("checking output: %w", err)
	}

	if len(c.Packages) == 0 {
		return fmt.Errorf("at least one package must be specified")
	}

	seen := make(map[string]bool)

	for _, p := range c.Packages {
		if _, ok := seen[p.Path]; ok {
			return fmt.Errorf("package %q is specified more than once", p.Path)
		}

		if err := p.check(); err != nil {
			return fmt.Errorf("checking package %q: %w", p.Path, err)
		}

		seen[p.Path] = true
	}

	return nil
}

func (c Config) Paths() []string {
	paths := make([]string, 0, len(c.Packages))

	for _, p := range c.Packages {
		paths = append(paths, p.Path)
	}

	return paths
}

type Output struct {
	Dir      string `koanf:"dir"`
	Package  string `koanf:"package"`
	Filename string `koanf:"filename"`
	Naming   Naming `koanf:"naming"`
}

func (o Output) check() error {
	if strings.TrimSpace(o.Dir) == "" {
		return errors.New("output directory is required")
	}

	if strings.TrimSpace(o.Package) == "" {
		return errors.New("output package is required")
	}

	if !token.IsIdentifier(o.Package) {
		return errors.New("output package name must be a valid identifier")
	}

	if strings.TrimSpace(o.Filename) == "" {
		return errors.New("output filename is required")
	}

	if !strings.Contains(o.Filename, "{package}") {
		return errors.New("output filename must contain {package} placeholder")
	}

	if err := o.Naming.check(); err != nil {
		return fmt.Errorf("checking naming: %w", err)
	}

	return nil
}

type Naming struct {
	Suffix string `koanf:"suffix"`
}

func (n Naming) check() error {
	if strings.TrimSpace(n.Suffix) == "" {
		return errors.New("suffix is required")
	}

	for i, r := range n.Suffix {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return fmt.Errorf("suffix contains invalid character '%c' at position %d", r, i)
		}
	}

	return nil
}

type Package struct {
	Path string `koanf:"path"`
}

func (p Package) check() error {
	if strings.TrimSpace(p.Path) == "" {
		return errors.New("package path is required")
	}

	return nil
}
