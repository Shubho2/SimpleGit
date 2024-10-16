package executing

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/codecrafters-io/git-starter-go/cmd/util/gitpath"
)

// ReadTreeOrBlobObject: reads the tree object from the .git/objects directory.
func ReadTreeOrBlobObject(shaDigest string) ([]byte, error) {
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
func WriteBlobObject(fileName string) ([]byte, error) {
	slog.Info("Writing blob object of", "fileName", fileName)

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		slog.Error("Error reading ", "file", fileName)
		return nil, err
	}
	
	contentToWrite := fmt.Sprintf("blob %d\x00%s", len(bytes), string(bytes))
	slog.Info("Content to write in blob object file", "contentToWrite", contentToWrite)

	shaDigest, err := write([]byte(contentToWrite))
	if err != nil {
		slog.Error("Error writing blob object of", "file", fileName)
		return nil, err
	}

	slog.Info("Successfully wrote blob object")
	return shaDigest, nil
}

func WriteTreeObject(root string) ([]byte, error) {
	slog.Info("Writing tree object of", "root", root)
	
	files, err := os.ReadDir(root)
	if err != nil {
		slog.Error("Error reading ", "directory", root)
		return nil, err
	}
	
	var bytes []byte
	for _, file := range files {
		if file.Name() == gitpath.Git {
			continue
		}
		
		var shaHash []byte
		if file.IsDir() {
			shaHash, err = WriteTreeObject(strings.Join([]string{root, file.Name()}, "/"))
			if err != nil {
				slog.Error("Error writing ", "file", file.Name())
				return nil, err
			}
		} else {
			shaHash, err = WriteBlobObject(strings.Join([]string{root, file.Name()}, "/"))
			if err != nil {
				slog.Error("Error writing ", "file", file.Name())
				return nil, err
			}
		}

		b2 := append([]byte(fmt.Sprintf("%s %s\x00", getFileMode(file), file.Name())), shaHash...)
		bytes = append(bytes, b2...)
	}
	
	dataToWrite := append([]byte(fmt.Sprintf("tree %d\x00", len(bytes))), bytes...)
	slog.Info("Content to write in tree object file", "contentToWrite", string(dataToWrite))

	shaDigest, err := write(dataToWrite)
	if err != nil {
		slog.Error("Error writing tree object of", "root", root)
		return nil, err
	}

	slog.Info("Successfully wrote tree object")
	return shaDigest, nil
}

func CommitTreeObject(treeShaDigest string, parentShaDigest string, message string) ([]byte, error) {
	slog.Info("Committing tree object of", "treeShaDigest", treeShaDigest)
	
	currentTime := time.Now().Unix()
	timezone, _ := time.Now().Local().Zone()
	author := "test"
	authorEmail := "test@gmail.com"

	commitData := fmt.Sprintf("tree %s\nparent %s\nauthor %s <%s> %s %s\ncommitter %s <%s> %s %s\n\n%s\n",
										treeShaDigest,
										parentShaDigest, 
										author, 
										authorEmail, 
										fmt.Sprint(currentTime), 
										timezone,
										author,
										authorEmail,
										fmt.Sprint(currentTime),
										timezone,
										message)
	content := []byte(fmt.Sprintf("commit %d\x00", len(commitData)))
	content = append(content, []byte(commitData)...)

	slog.Info("Content to write in commit object file", "contentToWrite", string(content))
	shaDigest, err := write(content)
	if err != nil {
		slog.Error("Error committing tree object of", "treeShaDigest", treeShaDigest)
		return nil, err
	}
	slog.Info("Successfully committed tree object")
	return shaDigest, nil
}

//******** Private Functions ********//

func write(bytes []byte) ([]byte, error) {
	shaDigest := calculateShaDigest(bytes)
	hexDigest := hex.EncodeToString(shaDigest)
	slog.Info("Calculated sha digest", "shaDigest", hexDigest)
	folderPath, err := createPathFrom(hexDigest)
	if err != nil {
		slog.Error("Error creating path from", "shaDigest", hexDigest)
		return nil, err
	}

	objectPath := strings.Join([]string{folderPath, hexDigest[2:]}, "/");
	file, err := os.OpenFile(objectPath, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0644);
	if err != nil {
		slog.Error("Error opening ", "file", objectPath)
		return nil, err
	}
	defer file.Close()

	writer := zlib.NewWriter(io.Writer(file))
	defer writer.Close()
	_, err = writer.Write(bytes)
	if err != nil {
		slog.Error("Error writing bytes to", "file", objectPath)
		return nil, err
	}

	slog.Info("Successfully wrote object to", "objectPath", objectPath)
	return shaDigest, nil
}

func getFileMode(file fs.DirEntry) string {
	if file.IsDir() {
		return "40000"
	} else {
		return "100644"
	}
}

func calculateShaDigest(bytes []byte) []byte {
	sha1 := sha1.New()
	sha1.Write(bytes)
	return sha1.Sum(nil)
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