package hash

import (
	"fmt"
	"log/slog"

	"github.com/codecrafters-io/git-starter-go/cmd/executing"
)

type HashObject struct {
	FileName string
}

func (ho HashObject) Execute(options map[string]bool) error {
	slog.Info("Called: HashObject.Execute()")
	shaDigest, err := executing.WriteBlobObject(ho.FileName)
	
	if err != nil {
		slog.Error("Error writing blob object", "err", err)
		return err
	}

	fmt.Print(shaDigest)
	return nil
}