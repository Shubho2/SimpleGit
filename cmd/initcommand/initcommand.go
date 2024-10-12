package initcommand

import (
	"log/slog"
	"os"

	"github.com/codecrafters-io/git-starter-go/cmd/util/gitpath"
)

// Init is a struct that implements the Executor interface
type Init struct {}

func (i Init) Execute(options map[string]bool) error {
	slog.Info("Called: Init.Execute()")
	paths := []string{gitpath.Git, gitpath.Objects, gitpath.Refs}

	for _, dir := range paths {
		if err := os.MkdirAll(dir, 0755); err != nil {
			slog.Error("Error creating directory:", "err", err)
			return err
		}
	}
	
	if err := os.WriteFile(gitpath.HEAD, getHeadFileContent(), 0644); err != nil {
		slog.Error("Error writing file:", "err", err)
		return err
	}
	
	slog.Info("Initialized git directory")
	return nil
}

func getHeadFileContent() []byte {
	return []byte("ref: refs/heads/main\n");
}