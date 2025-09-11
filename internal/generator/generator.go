package generator

import (
	"fmt"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/nuvrel/moldable/cmd/moldable/app"
	"github.com/nuvrel/moldable/internal/astfile"
	"github.com/nuvrel/moldable/internal/importset"
	"github.com/nuvrel/moldable/internal/pkgload"
	"github.com/nuvrel/moldable/internal/reporter"
	"github.com/nuvrel/moldable/internal/structcollector"
	"github.com/nuvrel/moldable/internal/typeast"
)

type Generator struct {
	config    app.Config
	reporter  reporter.Reporter
	collector *structcollector.StructCollector
	loader    *pkgload.Loader
	writer    *astfile.Writer
}

func New(cfg app.Config, rep reporter.Reporter) *Generator {
	return &Generator{
		config:    cfg,
		reporter:  rep,
		collector: structcollector.New(),
		loader:    pkgload.NewLoader(),
		writer:    astfile.NewWriter(),
	}
}

func (g Generator) Generate() error {
	paths := g.config.Paths()

	if err := g.loader.Load(paths); err != nil {
		return fmt.Errorf("loading packages: %w", err)
	}

	for _, p := range g.config.Packages {
		pkg, err := g.loader.Package(p.Path)
		if err != nil {
			return fmt.Errorf("using package after loading: %w", err)
		}

		if err := g.processPackage(pkg); err != nil {
			return fmt.Errorf("processing package %q: %w", p.Path, err)
		}
	}

	return nil
}

func (g Generator) processPackage(pkg *types.Package) error {
	g.reporter.ProcessingPackage(pkg.Path())

	builder := astfile.New(g.config.Output.Package)

	is := importset.New()
	is.Import(pkg)

	structs := g.collector.Collect(pkg)

	// TODO(calmondev): maybe we can move this counting to the collector?
	generated := 0

	for _, ss := range structs {
		if !ss.HasMethods() {
			g.reporter.SkippedStruct(ss.TypeName.Name(), "no methods")

			continue
		}

		typeast.TraverseTypeParams(ss.TypeParams, is.Import)
		typeast.TraverseFuncs(ss.Methods, is.Import)

		name := ss.TypeName.Name() + g.config.Output.Naming.Suffix

		builder.AddInterface(&astfile.InterfaceSpec{
			Name:    name,
			Methods: ss.Methods,
		})

		g.reporter.GeneratedInterface(name, ss.TypeName.Name(), len(ss.Methods))

		generated++
	}

	for path, alias := range is.Imports() {
		builder.AddImport(&astfile.ImportSpec{
			Path:  path,
			Alias: alias,
		})
	}

	if !builder.HasInterfaces() {
		return nil
	}

	file, err := builder.Build(is.Qualifier)
	if err != nil {
		return fmt.Errorf("building ast file: %w", err)
	}

	path := filepath.Join(
		g.config.Output.Dir,
		strings.ReplaceAll(g.config.Output.Filename, "{package}", pkg.Name()),
	)

	if err := g.writer.Write(file, path); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	g.reporter.PackageCompleted(pkg.Path(), generated)

	return nil
}
