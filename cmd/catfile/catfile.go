package catfile

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/codecrafters-io/git-starter-go/cmd/executing"
)

type CatFile struct {
	BlobShaDigest string
}

func (cf CatFile) Execute(options map[string]bool) error {
	slog.Info("Called: CatFile.Execute()")
	bytes, err := executing.ReadTreeObject(cf.BlobShaDigest)
	if err != nil {
		return err
	}

	index := strings.Index(string(bytes), "\x00")
	fmt.Print(string(bytes[index+1:]))
	return nil
}