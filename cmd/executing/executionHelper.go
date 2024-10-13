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
	file.Close()
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
	
	contentToWrite := fmt.Sprintf("blob %d\x00%s", len(bytes), string(bytes))
	slog.Info("Content to write", "contentToWrite", contentToWrite)

	shaDigest, err := write([]byte(contentToWrite))
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
	slog.Info("Calculated sha digest", "shaDigest", shaDigest)
	
	folderPath, err := createPathFrom(shaDigest)
	if err != nil {
		slog.Error("Error creating path from", "shaDigest", shaDigest)
		return "", err
	}

	objectPath := strings.Join([]string{folderPath, shaDigest[2:]}, "/");
	file, err := os.OpenFile(objectPath, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0644);

	if err != nil {
		slog.Error("Error opening ", "file", objectPath)
		return "", err
	}

	writer := zlib.NewWriter(io.Writer(file));
	_, err = writer.Write(bytes);
	if err != nil {
		slog.Error("Error writing bytes to", "file", objectPath)
		return "", err
	}

	writer.Close();
	file.Close();
	return shaDigest, nil
}


func calculateShaDigest(bytes []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(bytes));
}

func createPathFrom(shaDigest string) (string, error) {
	var folderPath string = strings.Join([]string{gitpath.Objects, shaDigest[0:2]}, "/");
	err := os.MkdirAll(folderPath, 0755);
	if err != nil {
		slog.Error("Error creating directory", "objectPath", folderPath);
		return "", err
	}

	slog.Info("Created directory at", "objectPath", folderPath);
	return folderPath, nil
}

// Gets the path to the file URL from the SHA digest.
func getPathToFileURLFrom(shaDigest string) string {
	slog.Info("Decoding path from digest", "shaDigest", shaDigest)
	return strings.Join([]string{gitpath.Objects, shaDigest[0:2], shaDigest[2:]}, "/")
}