package main

import (
	"fmt"
	"github.com/jonasrdl/bookmark-sync/internal/browser/chromium"
	"github.com/jonasrdl/bookmark-sync/internal/browser/firefox"
)

func main() {
	chromium := chromium.ChromiumBrowser{}
	chromiumFilepath, _ := chromium.GetBookmarksFilepath()
	fmt.Println("chromium:", chromiumFilepath)

	firefox := firefox.FirefoxBrowser{}
	firefoxFilepath, _ := firefox.GetBookmarksFilepath()
	fmt.Println("firefox:", firefoxFilepath)
}
