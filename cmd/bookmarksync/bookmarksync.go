package main

import (
	"fmt"
	"github.com/jonasrdl/bookmark-sync/internal"
	"github.com/jonasrdl/bookmark-sync/internal/browser/chromium"
	"github.com/jonasrdl/bookmark-sync/internal/browser/firefox"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "bookmark-sync",
	Short: "A tool to sync bookmarks between browsers",
	Run: func(cmd *cobra.Command, args []string) {
		// show help msg if no subcommand is provided
		cmd.Help()
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync bookmarks between browsers",
	Run: func(cmd *cobra.Command, args []string) {
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")

		if from == "" || to == "" {
			log.Fatal("Both source and destination browsers must be specified")
		}

		syncBookmarks(from, to)
	},
}

func init() {
	syncCmd.Flags().String("from", "", "Source browser (chromium or firefox)")
	syncCmd.Flags().String("to", "", "Destination browser (chromium or firefox)")

	rootCmd.AddCommand(syncCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func syncBookmarks(from, to string) {
	var sourceBrowser, destBrowser internal.Browser

	fmt.Println("PING 1")

	switch from {
	case "chromium":
		sourceBrowser = &chromium.ChromiumBrowser{}
	case "firefox":
		sourceBrowser = &firefox.FirefoxBrowser{}
	default:
		log.Fatal("Invalid source browser specified")
	}

	switch to {
	case "chromium":
		destBrowser = &chromium.ChromiumBrowser{}
	case "firefox":
		destBrowser = &firefox.FirefoxBrowser{}
	default:
		log.Fatal("Invalid destination browser specified")
	}

	fmt.Println("PING 2")

	sourceBrowserFilepath, _ := sourceBrowser.GetBookmarksFilepath()
	sourceBookmarks, err := sourceBrowser.ParseJSON(sourceBrowserFilepath)
	if err != nil {
		log.Fatal("error parsing source browser bookmarks:", err)
	}

	fmt.Println("PING 3")

	destBrowserFilepath, _ := destBrowser.GetBookmarksFilepath()
	destBookmarks, err := destBrowser.ParseJSON(destBrowserFilepath)
	if err != nil {
		log.Fatal("error parsing destination browser bookmarks:", err)
	}

	mergedBookmarks := internal.MergeBookmarks(sourceBookmarks, destBookmarks)

	fmt.Println("PING 4")

	if err := destBrowser.UpdateJSON(mergedBookmarks); err != nil {
		log.Fatal("Error updating destination browser bookmarks", err)
	}
	fmt.Println("Bookmarks synced successfully")
}
