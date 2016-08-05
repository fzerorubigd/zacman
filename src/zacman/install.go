package main

import (
	"flagx"

	"plugin"

	"github.com/ogier/pflag"
)

var (
	installTarget string
	installUpdate bool
)

func installPackage() error {
	root := makeRoot()
	err := plugin.InstallPlugin(root, installTarget, installUpdate)
	if err == nil {
		err = plugin.MkZshRC(root)
	}
	return err
}

func init() {
	f := pflag.NewFlagSet("install", pflag.ExitOnError)
	f.StringVarP(&installTarget, "package", "p", "", "the package name to install")
	f.BoolVarP(&installUpdate, "update", "u", false, "Update if already installed?")
	flagx.RegisterCommand("install", f, installPackage)
}
