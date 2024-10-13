package lstree

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/codecrafters-io/git-starter-go/cmd/executing"
)

type LsTree struct {
	TreeShaDigest string
}

func (lt LsTree) Execute(options map[string]bool) error {
	slog.Info("Called: LsTree.Execute()")
	bytes, err := executing.ReadTreeOrBlobObject(lt.TreeShaDigest)
	if err != nil {
		return err
	}

	index := strings.Index(string(bytes), "\x00")
	data := string(bytes[index+1:])
	
	for len(data) > 0 {
		output := parseGitTreeEntry(data)
		data = output[0]
		fmt.Println(output[3])
	}
	return nil
}

func parseGitTreeEntry(data string) []string {
	index := strings.Index(data, " ")
	mode := data[0:index]
	data = data[index+1:]
	index = strings.Index(data, "\x00")
	name := data[0:index]
	data = data[index+1:]
	shaDigest := data[0:20];
	shaDigest = fmt.Sprintf("%x", shaDigest)
	
	objectType := "blob"
	if mode[0] == '4' {
		objectType = "tree"
	}
	data = data[20:]
	return []string{data, mode, objectType, name, shaDigest};
}