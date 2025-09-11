package app

import "embed"

//go:embed templates/*.tmpl
var Templates embed.FS

const (
	Name               = "moldable"
	ConfigFile         = Name + ".yaml"
	ConfigFileTemplate = ConfigFile + ".tmpl"
)

const (
	ConfigFileFlag = "config-file"
	ForceFlag      = "force"
)
