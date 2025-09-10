package reporter

import "github.com/charmbracelet/log"

type Log struct {
	logger *log.Logger
}

var (
	_ Reporter = (*Log)(nil)
)

func NewLog(logger *log.Logger) *Log {
	return &Log{
		logger: logger,
	}
}

func (l Log) ProcessingPackage(name string) {
	l.logger.Info("processing package", "path", name)
}

func (l Log) GeneratedInterface(interfaceName, structName string, methodCount int) {
	l.logger.Info("generated interface", "name", interfaceName, "from_struct", structName, "method_count", methodCount)
}

func (l Log) SkippedStruct(structName, reason string) {
	l.logger.Info("skipped struct", "name", structName, "reason", reason)
}

func (l Log) PackageCompleted(packagePath string, interfaceCount int) {
	l.logger.Info("completed package", "path", packagePath, "interface_count", interfaceCount)
}
