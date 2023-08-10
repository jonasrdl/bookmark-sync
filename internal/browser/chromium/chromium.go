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

func (c *ChromiumBrowser) UpdateJSON(bookmarks []internal.Bookmark) error {
	bookmarksFilePath, _ := c.GetBookmarksFilepath()

	// Read existing bookmarks from the Chromium bookmarks file
	existingBookmarks, err := c.ParseJSON(bookmarksFilePath)
	if err != nil {
		log.Printf("error parsing existing bookmarks: %v\n", err)
		return err
	}

	// Merge the existing Chromium bookmarks with the new ones
	mergedBookmarks := mergeBookmarks(existingBookmarks, bookmarks)

	// Generate the Chromium bookmarks JSON format
	chromiumBookmarks := struct {
		Checksum string `json:"checksum"`
		Roots    map[string]struct {
			Children []internal.Bookmark `json:"children"`
		} `json:"roots"`
		Version int `json:"version"`
	}{
		Checksum: "",
		Roots: map[string]struct {
			Children []internal.Bookmark `json:"children"`
		}{
			"bookmark_bar": {
				Children: mergedBookmarks,
			},
		},
		Version: 1,
	}

	// Serialize bookmarks to JSON
	jsonData, err := json.MarshalIndent(chromiumBookmarks, "", "  ")
	if err != nil {
		log.Printf("error marshaling bookmarks to JSON: %v\n", err)
		return err
	}

	// Write JSON data to the bookmarks file
	err = os.WriteFile(bookmarksFilePath, jsonData, 0644)
	if err != nil {
		log.Printf("error writing JSON data to file: %v\n", err)
		return err
	}

	return nil
}

func mergeBookmarks(existing, new []internal.Bookmark) []internal.Bookmark {
	// Create a map to keep track of existing bookmarks using their GUID
	existingMap := make(map[string]internal.Bookmark)
	for _, bookmark := range existing {
		existingMap[bookmark.ID] = bookmark
	}

	// Merge the new bookmarks into the existing bookmarks map
	for _, bookmark := range new {
		// Check if the bookmark already exists based on its GUID
		if _, exists := existingMap[bookmark.ID]; !exists {
			existingMap[bookmark.ID] = bookmark
		}
	}

	// Convert the merged map back into a slice of bookmarks
	var merged []internal.Bookmark
	for _, bookmark := range existingMap {
		merged = append(merged, bookmark)
	}

	return merged
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
