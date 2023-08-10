package internal

type Browser interface {
	GetBookmarksFilepath() (string, error)
}
