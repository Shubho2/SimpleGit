package executing

import (
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/codecrafters-io/git-starter-go/cmd/util/gitpath"
)

// ReadTreeObject: reads the tree object from the .git/objects directory.
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

// WriteBlobObject: writes the blob object to the .git/objects directory.
func WriteBlobObject(fileName string) (string, error) {
	slog.Info("Writing blob object using", "fileName", fileName)

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		slog.Error("Error reading ", "file", fileName)
		return "", err
	}

	contentToWrite := []byte("blob " + string(len(bytes)) + "\x00" + string(bytes)) 
	shaDigest, err := write(contentToWrite)
	if err != nil {
		slog.Error("Error writing ", "file", fileName)
		return "", err
	}

	slog.Info("Successfully wrote blob object")
	return shaDigest, nil
}


//******** Private Functions ********//

func write(bytes []byte) (string, error) {
	shaDigest := calculateShaDigest(bytes)
	pathToFileURL := createPathFrom(shaDigest);

	file, err := os.Open(pathToFileURL)
	if err != nil {
		slog.Error("Error opening ", "file", pathToFileURL)
		return "", err
	}

	w := zlib.NewWriter(io.Writer(file));
	_, err = w.Write(bytes);
	if err != nil {
		slog.Error("Error writing bytes to", "file", pathToFileURL)
		return "", err
	}

	return shaDigest, nil
}


func calculateShaDigest(bytes []byte) string {
	sha := sha1.New();
	sha.Write(bytes);
	return fmt.Sprintf("%x", sha.Sum(nil));
}

func createPathFrom(shaDigest string) string {
	var objectPath string = strings.Join([]string{gitpath.Objects, shaDigest[0:2]}, "/");
	os.MkdirAll(objectPath, 0755);
	slog.Info("Created directory", "objectPath", objectPath);
	return objectPath
}

// Gets the path to the file URL from the SHA digest.
func getPathToFileURLFrom(shaDigest string) string {
	slog.Info("Decoding path from digest", "shaDigest", shaDigest)
	return strings.Join([]string{gitpath.Objects, shaDigest[0:2], shaDigest[2:]}, "/")
}