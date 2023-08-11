package chromium

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/jonasrdl/bookmark-sync/internal"
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

	existingBookmarks, err := c.ParseJSON(bookmarksFilePath)
	if err != nil {
		log.Printf("error parsing existing bookmarks: %v\n", err)
		return err
	}

	mergedBookmarks := mergeBookmarks(existingBookmarks, bookmarks)

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

	jsonData, err := json.MarshalIndent(chromiumBookmarks, "", "  ")
	if err != nil {
		log.Printf("error marshaling bookmarks to JSON: %v\n", err)
		return err
	}

	err = os.WriteFile(bookmarksFilePath, jsonData, 0o644)
	if err != nil {
		log.Printf("error writing JSON data to file: %v\n", err)
		return err
	}

	return nil
}

func mergeBookmarks(existing, new []internal.Bookmark) []internal.Bookmark {
	existingMap := make(map[string]internal.Bookmark)
	for _, bookmark := range existing {
		existingMap[bookmark.ID] = bookmark
	}

	for _, bookmark := range new {
		if _, exists := existingMap[bookmark.ID]; !exists {
			existingMap[bookmark.ID] = bookmark
		}
	}

	var merged []internal.Bookmark
	for _, bookmark := range existingMap {
		merged = append(merged, bookmark)
	}

	return merged
}
