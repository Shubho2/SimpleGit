package main

import (
	"log/slog"
	"os"

	"github.com/codecrafters-io/git-starter-go/cmd/catfile"
	"github.com/codecrafters-io/git-starter-go/cmd/commit"
	"github.com/codecrafters-io/git-starter-go/cmd/executing"
	"github.com/codecrafters-io/git-starter-go/cmd/hash"
	"github.com/codecrafters-io/git-starter-go/cmd/initcommand"
	"github.com/codecrafters-io/git-starter-go/cmd/lstree"
	"github.com/codecrafters-io/git-starter-go/cmd/util/command"
	"github.com/codecrafters-io/git-starter-go/cmd/write"
)

func main() {
	slog.Debug("Called: main()")

	if len(os.Args) < 2 {
		slog.Info("usage: mygit <command> [<args>...]")
		os.Exit(1)
	}

	var ex executing.Executor
	options := make(map[string]bool)

	switch _command := os.Args[1]; _command {
	case command.Init:
		ex = initcommand.Init{}
		if err := commandExecutor(ex, options); err != nil {
			slog.Error("Error executing init command", "err", err)
			os.Exit(1)
		}
	case command.CatFile:
		if(os.Args[2] == "-p") { 
			options["pretty_print"] = true
			ex = catfile.CatFile{BlobShaDigest: os.Args[3]}
		} else {
			ex = catfile.CatFile{BlobShaDigest: os.Args[2]}
		}

		if err := commandExecutor(ex, options); err != nil {
			slog.Error("Error executing cat-file command", "err", err)
			os.Exit(1)
		}
	case command.HashObject:
		if(os.Args[2] == "-w") {
			options["write"] = true
			ex = hash.HashObject{FileName: os.Args[3]}
		} else {
			ex = hash.HashObject{FileName: os.Args[2]}
		}

		if err := commandExecutor(ex, options); err != nil {
			slog.Error("Error executing hash-object command", "err", err)
			os.Exit(1)
		}
	case command.LsTree:
		if(os.Args[2] == "--name-only") { 
			options["name_only"] = true
			ex = lstree.LsTree{TreeShaDigest: os.Args[3]}
		} else {
			ex = lstree.LsTree{TreeShaDigest: os.Args[2]}
		}

		if err := commandExecutor(ex, options); err != nil {
			slog.Error("Error executing ls-tree command", "err", err)
			os.Exit(1)
		}
	case command.WriteTree:
		ex = write.WriteTree{}
		if err := commandExecutor(ex, options); err != nil {
			slog.Error("Error executing write-tree command", "err", err)
			os.Exit(1)
		}
	case command.CommitTree:
		ex = commit.CommitTree { 	
						TreeShaDigest: os.Args[2], 
						ParentShaDigest: os.Args[4], 
						Message: os.Args[6],
					}
		if err := commandExecutor(ex, options); err != nil {
			slog.Error("Error executing commit-tree command", "err", err)
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
