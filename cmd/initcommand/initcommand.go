package initcommand

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/cmd/util/gitpath"
)

type InitType struct {

}

func (i InitType) Execute(options map[string]bool) error {

	paths := []string{gitpath.Git, gitpath.Objects, gitpath.Refs}

	for _, dir := range paths {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			return err
		}
	}
	
	if err := os.WriteFile(gitpath.HEAD, getHeadFileContent(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		return err
	}
	
	fmt.Println("Initialized git directory")
	return nil
}

func getHeadFileContent() []byte {
	return []byte("ref: refs/heads/main\n");
}