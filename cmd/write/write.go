package write

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/codecrafters-io/git-starter-go/cmd/executing"
)

type WriteTree struct {}

func (w WriteTree) Execute(options map[string]bool) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	
	hash, err := executing.WriteTreeObject(currentDir)
	if err != nil {
		return err
	}
	
	fmt.Print(hex.EncodeToString(hash))
	return nil
}