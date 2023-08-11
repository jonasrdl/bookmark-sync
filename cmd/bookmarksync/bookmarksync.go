package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/jonasrdl/bookmark-sync/internal"
	"github.com/jonasrdl/bookmark-sync/internal/browser/chromium"
	"github.com/jonasrdl/bookmark-sync/internal/browser/firefox"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bookmark-sync",
	Short: "A tool to sync bookmarks between browsers",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync bookmarks between browsers",
	Run: func(cmd *cobra.Command, args []string) {
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")

		//destProfile, _ := cmd.Flags().GetString("dest-profile")

		if from == "" || to == "" {
			log.Fatal("Both source and destination browsers must be specified")
		}

		syncBookmarks(from, to)
	},
}

var listProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: "List available browser profiles",
	Run: func(cmd *cobra.Command, args []string) {
		browser, _ := cmd.Flags().GetString("browser")

		switch browser {
		case "chromium":
			listChromiumProfiles()
		case "firefox":
			listFirefoxProfiles()
		default:
			log.Fatal("Invalid browser specified")
		}
	},
}

func init() {
	syncCmd.Flags().String("from", "", "Source browser (chromium or firefox)")
	syncCmd.Flags().String("to", "", "Destination browser (chromium or firefox)")
	syncCmd.Flags().String("source-profile", "", "Source browser profile")
	syncCmd.Flags().String("dest-profile", "", "Destination browser profile")

	listProfilesCmd.Flags().String("browser", "", "Browser to list profiles (chromium or firefox)")

	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(listProfilesCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func syncBookmarks(from, to string) {
	var sourceBrowser, destBrowser internal.Browser

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

	sourceBrowserFilepath, _ := sourceBrowser.GetBookmarksFilepath()
	sourceBookmarks, err := sourceBrowser.ParseJSON(sourceBrowserFilepath)
	if err != nil {
		log.Fatal("error parsing source browser bookmarks:", err)
	}

	destBrowserFilepath, _ := destBrowser.GetBookmarksFilepath()
	destBookmarks, err := destBrowser.ParseJSON(destBrowserFilepath)
	if err != nil {
		log.Fatal("error parsing destination browser bookmarks:", err)
	}

	mergedBookmarks := internal.MergeBookmarks(sourceBookmarks, destBookmarks)

	if err := destBrowser.UpdateJSON(mergedBookmarks); err != nil {
		log.Fatal("Error updating destination browser bookmarks", err)
	}
	fmt.Println("Bookmarks synced successfully")
}

func listChromiumProfiles() {
	fmt.Println("Currently not supported")
}

func listFirefoxProfiles() {
	fmt.Println("List of available Firefox profiles:")
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	firefoxDir := filepath.Join(usr.HomeDir, ".mozilla", "firefox")
	profilesIniPath := filepath.Join(firefoxDir, "profiles.ini")

	profiles, err := readFirefoxProfilesIni(profilesIniPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, profile := range profiles {
		fmt.Println(profile)
	}
}

func readFirefoxProfilesIni(profilesIniPath string) ([]string, error) {
	file, err := os.Open(profilesIniPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var profiles []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Path=") {
			profilePath := strings.TrimPrefix(line, "Path=")
			profiles = append(profiles, profilePath)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}
