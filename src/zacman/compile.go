package main

import (
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func compileEntry(cmd *cobra.Command, args []string) {
	makeRoot()
	p, err := loadSnapShot("master")
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}

	if concurrent < 1 {
		concurrent = 1
	}
	//TODO : better handle this
	compileSnapshot(p, concurrent)
	ioutil.WriteFile(root+"/zacman.zsh", []byte(buildLoadScipt(p)), 0644)

	logrus.Info("compile is done, check the log if there is a problem")
}

func initCompileCommand() *cobra.Command {
	var compile = &cobra.Command{
		Use:   "compile",
		Short: "try to compile all plugins",
		Long:  `try to compile all plugins`,
		Run:   compileEntry,
	}

	compile.Flags().UintVarP(&concurrent, "concurrent", "c", 5, "how many clone in parallel?")

	return compile
}
