package executing

import (
	"compress/zlib"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/codecrafters-io/git-starter-go/cmd/util/gitpath"
)

// ReadTreeObject: reads the tree object from the .git/objects directory.
// shaDigest: The SHA digest.
// Returns: The tree object as a bytes.Buffer.
func ReadTreeObject(shaDigest string) ([]byte, error) {
	slog.Info("Reading tree object using", "shaDigest", shaDigest)
	pathToFileURL := getPathToFileURLFrom(shaDigest)

	file, err := os.Open(pathToFileURL)
	if err != nil {
		slog.Error("Error opening ", "file", pathToFileURL)
		return nil, err
	}

	zlibReader, err := zlib.NewReader(io.Reader(file))
	if err != nil {
		slog.Error("Error decompressing ", "file", pathToFileURL)
		return nil, err
	}

	buffer, err := io.ReadAll(zlibReader)
	if err != nil {
		slog.Error("Error reading ", "file", pathToFileURL)
		return nil, err
	}

	slog.Info("Successfully read tree object")
	zlibReader.Close()
	return buffer, nil
}


//******** Private Functions ********//

/**
* Gets the path to the file URL from the SHA digest.
* shaDigest: The SHA digest.
*/
func getPathToFileURLFrom(shaDigest string) string {
	slog.Info("Decoding path from digest", "shaDigest", shaDigest)
	return strings.Join([]string{gitpath.Objects, shaDigest[0:2], shaDigest[2:]}, "/")
}