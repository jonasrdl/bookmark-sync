package internal

type Browser interface {
	GetBookmarksFilepath() (string, error)
	ParseJSON(path string) ([]Bookmark, error)
	UpdateJSON(bookmarks []Bookmark) error
}
