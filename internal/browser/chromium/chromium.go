package chromium

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
)

type ChromiumBrowser struct {
	UserProfileDir string
}

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
