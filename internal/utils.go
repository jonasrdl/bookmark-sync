package internal

func MergeBookmarks(existing, new []Bookmark) []Bookmark {
	existingMap := make(map[string]Bookmark)

	for _, bookmark := range existing {
		existingMap[bookmark.ID] = bookmark
	}

	for _, bookmark := range new {
		if _, exists := existingMap[bookmark.ID]; !exists {
			existingMap[bookmark.ID] = bookmark
		}
	}

	var merged []Bookmark
	for _, bookmark := range existingMap {
		merged = append(merged, bookmark)
	}

	return merged
}
