package main

import (
	"flagx"
	"plugin"

	"github.com/ogier/pflag"
)

var (
	removeTarget string
)

func removePlugin() error {
	return plugin.RemovePlugin(makeRoot(), removeTarget)
}

func init() {
	f := pflag.NewFlagSet("remove", pflag.ExitOnError)
	f.StringVarP(&removeTarget, "package", "p", "", "the package name to remove")
	flagx.RegisterCommand("remove", f, removePlugin)
}
