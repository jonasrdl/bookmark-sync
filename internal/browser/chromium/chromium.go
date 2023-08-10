package chromium

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/jonasrdl/bookmark-sync/internal"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type ChromiumBrowser struct{}

func (c *ChromiumBrowser) GetBookmarksFilepath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	profileDir := filepath.Join(usr.HomeDir, ".config", "chromium", "Default")
	bookmarksFilePath := filepath.Join(profileDir, "Bookmarks")
	if _, err := os.Stat(bookmarksFilePath); err != nil {
		return "", errors.New("chromium bookmarks file not found")
	}
	return bookmarksFilePath, nil
}

func (c *ChromiumBrowser) ParseJSON(path string) ([]internal.Bookmark, error) {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		log.Printf("error reading json file: %v\n", err)
		return nil, err
	}

	var data struct {
		Roots struct {
			BookmarkBar struct {
				Children []internal.Bookmark `json:"children"`
			} `json:"bookmark_bar"`
		} `json:"roots"`
	}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Printf("error unmarshaling json: %v\n", err)
		return nil, err
	}

	return data.Roots.BookmarkBar.Children, nil
}

// CalculateChecksum calculates the SHA-256 checksum for the content of a file.
// It reads the contents of the file at the specified filePath, computes the hash,
// and returns the checksum as a hexadecimal string.
// If an error occurs while reading the file or computing the hash, an error is returned.
func CalculateChecksum(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}
