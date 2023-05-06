package cmd

import (
	"context"
	"flag"
	"github.com/google/subcommands"
)

type verboseKey struct{}

var VerboseKey = verboseKey{}

// Execute sets up the command chain and runs it.
func Execute(ctx context.Context) subcommands.ExitStatus {
	for _, command := range [...]subcommands.Command{
		subcommands.CommandsCommand(),
		subcommands.FlagsCommand(),
		subcommands.HelpCommand(),
		&ImageCommand{},
		&ProgramCommand{},
	} {
		subcommands.Register(command, "")
	}

	verbose := flag.Bool("v", true, "Enable verbose mode")
	subcommands.ImportantFlag("v")
	flag.Parse()
	ctx = context.WithValue(ctx, VerboseKey, *verbose)

	return subcommands.Execute(ctx)
}
