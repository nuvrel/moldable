package reporter

type Reporter interface {
	ProcessingPackage(packagePath string)
	GeneratedInterface(interfaceName, structName string, methodCount int)
	SkippedStruct(structName, reason string)
	PackageCompleted(packagePath string, interfaceCount int)
}
