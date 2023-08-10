package chromium

import (
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
