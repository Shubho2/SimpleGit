package commit

import (
	"encoding/hex"
	"fmt"

	"github.com/codecrafters-io/git-starter-go/cmd/executing"
)

type CommitTree struct {
	TreeShaDigest string
	ParentShaDigest string
	Message string
}

func (ct CommitTree) Execute(options map[string]bool) error {
	hash, err := executing.CommitTreeObject(ct.TreeShaDigest, ct.ParentShaDigest, ct.Message)
	if err != nil {
		return err
	}
	
	fmt.Print(hex.EncodeToString(hash))
	return nil
}