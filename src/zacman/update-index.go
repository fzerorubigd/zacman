package main

import (
	"flagx"
	"plugin"
)

func updateIndex() error {
	return plugin.UpdateIndex(makeRoot())
}

func init() {
	//f := pflag.NewFlagSet("update-index", pflag.ExitOnError)
	flagx.RegisterCommand("update-index", nil, updateIndex)
}
