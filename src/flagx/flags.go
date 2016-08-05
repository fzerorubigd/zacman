package flagx

import (
	"fmt"

	"assert"

	"os"

	"github.com/ogier/pflag"
)

var (
	commands = make(map[string]*cmd)
)

type Runnable func() error

type cmd struct {
	Flags *pflag.FlagSet
	Entry Runnable
}

func RegisterCommand(sub string, f *pflag.FlagSet, r Runnable) {
	e := &cmd{
		Flags: f,
		Entry: r,
	}
	if e.Flags == nil {
		e.Flags = pflag.NewFlagSet(sub, pflag.ExitOnError)
	}

	assert.NotNil(e.Entry)
	commands[sub] = e
}

func RunCommand(sub string, args ...string) error {
	e, ok := commands[sub]
	if !ok {
		return fmt.Errorf("there is no such command : %s", sub)
	}
	err := e.Flags.Parse(args)
	if err != nil {
		return err
	}
	
	return e.Entry()
}

func PrintUsage() {
	panic("TODO")
	os.Exit(2)
}

func Parse() error {
	args := os.Args
	if len(args) < 2 || args[1] == "-h" || args[1] == "--help" {
		PrintUsage()
	}
	post := []string{}
	if len(args) > 2 {
		post = append(post, args[2:]...)
	}
	return RunCommand(args[1], post...)
}
