package firefox

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/jonasrdl/bookmark-sync/internal"
	_ "github.com/mattn/go-sqlite3"
)

type FirefoxBrowser struct{}

func (f *FirefoxBrowser) ParseJSON(path string) ([]internal.Bookmark, error) {
	bookmarks, err := readBookmarksFromSQLite(path)
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
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

func (f *FirefoxBrowser) UpdateJSON(bookmarks []internal.Bookmark) error {
	bookmarksFilePath, _ := f.GetBookmarksFilepath()

	firefoxBm, err := readBookmarksFromSQLite(bookmarksFilePath)
	if err != nil {
		log.Printf("error reading bookmarks from SQLite: %v\n", err)
		return err
	}

	mergedBookmarks := mergeBookmarks(firefoxBm, bookmarks)

	firefoxBookmarks := generateFirefoxStyleJSON(mergedBookmarks)

	jsonData, err := json.MarshalIndent(firefoxBookmarks, "", "  ")
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

func generateFirefoxStyleJSON(bookmarks []internal.Bookmark) map[string]interface{} {
	root := make(map[string]interface{})
	root["checksum"] = ""
	roots := make(map[string]interface{})
	bookmarkBar := make(map[string]interface{})
	bookmarkBar["children"] = bookmarks
	bookmarkBar["date_added"] = ""
	bookmarkBar["date_modified"] = ""
	bookmarkBar["id"] = ""
	bookmarkBar["name"] = "bookmark_bar"
	bookmarkBar["type"] = "folder"
	roots["bookmark_bar"] = bookmarkBar
	other := make(map[string]interface{})
	other["children"] = []interface{}{}
	other["date_added"] = ""
	other["date_modified"] = ""
	other["id"] = ""
	other["name"] = "Other bookmarks"
	other["type"] = "folder"
	roots["other"] = other
	synced := make(map[string]interface{})
	synced["children"] = []interface{}{}
	synced["date_added"] = ""
	synced["date_modified"] = ""
	synced["id"] = ""
	synced["name"] = "Mobile bookmarks"
	synced["type"] = "folder"
	roots["synced"] = synced
	root["roots"] = roots
	root["version"] = 1

	return root
}

func readBookmarksFromSQLite(path string) ([]internal.Bookmark, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Printf("error opening SQLite database: %v\n", err)
		return nil, err
	}
	defer db.Close()

	var (
		rows     *sql.Rows
		attempts = 3
	)

	for i := 0; i < attempts; i++ {
		rows, err = db.Query("SELECT moz_bookmarks.id, moz_places.url, moz_bookmarks.title FROM moz_bookmarks INNER JOIN moz_places ON moz_bookmarks.fk = moz_places.id WHERE moz_bookmarks.type = 1")
		if err == nil {
			break
		}
		log.Printf("error querying bookmarks from database (attempt %d): %v\n", i+1, err)
		time.Sleep(time.Second) // Wait for a second before retrying
	}

	if rows == nil {
		return nil, err
	}

	defer rows.Close()

	var bookmarks []internal.Bookmark
	for rows.Next() {
		var id, url, title string
		err := rows.Scan(&id, &url, &title)
		if err != nil {
			log.Printf("error scanning bookmark rows: %v\n", err)
			return nil, err
		}

		bookmark := internal.Bookmark{
			ID:   id,
			URL:  url,
			Name: title,
		}
		bookmarks = append(bookmarks, bookmark)
	}

	return bookmarks, nil
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

	for _, path := range profilePaths {
		profilePath := filepath.Join(filepath.Dir(profilesIniPath), path)
		bookmarksFilePath := filepath.Join(profilePath, "places.sqlite")
		if _, err := os.Stat(bookmarksFilePath); err == nil {
			return profilePath, nil
		}
	}

	return "", errors.New("no valid profile with places.sqlite found")
}
