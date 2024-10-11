package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/cmd/executing"
	"github.com/codecrafters-io/git-starter-go/cmd/initcommand"
	"github.com/codecrafters-io/git-starter-go/cmd/util/command"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	// Interface
	var ex executing.Executor

	switch _command := os.Args[1]; _command {
	case command.Init:
		ex = initcommand.InitType{}
		err := ex.Execute(nil)
		if err != nil {
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", _command)
		os.Exit(1)
	}
}
