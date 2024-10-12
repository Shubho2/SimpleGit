package catfile

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/codecrafters-io/git-starter-go/cmd/executing"
)

type CatFile struct {
	blobShaDigest string
}

func (cf CatFile) Execute(options map[string]bool) error {
	slog.Debug("Called: CatFile.Execute()")
	bytes, err := executing.ReadTreeObject(cf.blobShaDigest)
	if err != nil {
		return err
	}

	index := strings.Index(string(bytes), "\x00")
	fmt.Print(bytes[index+1:])
	return nil
}