package firefox

import (
	"bufio"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

/**
TODO: sqlite has to be parsed, which is a pain in the ass
*/

type FirefoxBrowser struct {
	UserProfileDir string
}

func (f *FirefoxBrowser) GetBookmarksFilepath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	profileDir := filepath.Join(usr.HomeDir, ".mozilla", "firefox")

	profilesIniPath := filepath.Join(profileDir, "profiles.ini")
	profilePath, err := findValidProfilePath(profilesIniPath)
	if err != nil {
		return "", err
	}

	bookmarksFilePath := filepath.Join(profilePath, "places.sqlite")
	if _, err := os.Stat(bookmarksFilePath); err != nil {
		return "", errors.New("firefox bookmarks file not found")
	}
	return bookmarksFilePath, nil
}

func findValidProfilePath(profilesIniPath string) (string, error) {
	file, err := os.Open(profilesIniPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var defaultProfilePath string
	var profilePaths []string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Path=") {
			profilePath := strings.TrimPrefix(line, "Path=")
			profilePaths = append(profilePaths, profilePath)
			if defaultProfilePath == "" {
				defaultProfilePath = profilePath
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if defaultProfilePath == "" {
		return "", errors.New("default profile path not found in profiles.ini")
	}

	// Search for a valid places.sqlite file among the profiles
	for _, path := range profilePaths {
		profilePath := filepath.Join(filepath.Dir(profilesIniPath), path)
		bookmarksFilePath := filepath.Join(profilePath, "places.sqlite")
		if _, err := os.Stat(bookmarksFilePath); err == nil {
			return profilePath, nil
		}
	}

	return "", errors.New("no valid profile with places.sqlite found")
}
