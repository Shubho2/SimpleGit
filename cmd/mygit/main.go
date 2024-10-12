package main

import (
	"log/slog"
	"os"

	"github.com/codecrafters-io/git-starter-go/cmd/catfile"
	"github.com/codecrafters-io/git-starter-go/cmd/executing"
	"github.com/codecrafters-io/git-starter-go/cmd/initcommand"
	"github.com/codecrafters-io/git-starter-go/cmd/util/command"
)

func main() {
	slog.Debug("Called: main()")

	if len(os.Args) < 2 {
		slog.Info("usage: mygit <command> [<args>...]")
		os.Exit(1)
	}

	// Declaring an variable of Executor type
	var ex executing.Executor
	options := make(map[string]bool)

	switch _command := os.Args[1]; _command {
	case command.Init:
		ex = initcommand.Init{}
		if err := commandExecutor(ex, options); err != nil {
			os.Exit(1)
		}
	case command.CatFile:
		ex = catfile.CatFile{}
		if(os.Args[2] == "-p") { 
			options["pretty_print"] = true
		}
		if err := commandExecutor(ex, options); err != nil {
			os.Exit(1)
		}
	default:
		slog.Error("Unknown command", "command", _command)
		os.Exit(1)
	}
}

func commandExecutor(ex executing.Executor, options map[string]bool) error {
	slog.Debug("Called: commandExecutor()")
	return ex.Execute(options)
}
