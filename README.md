# Bookmark Sync
Bookmark Sync is a command-line tool that allows you to synchronize bookmarks between different web browsers, currently supporting Chromium and Firefox.

## Table of Contents
- [Installation](https://github.com/jonasrdl/bookmark-sync/tree/main#installation)
- [Usage](https://github.com/jonasrdl/bookmark-sync/tree/main#usage)
- [Command-Line Interface](https://github.com/jonasrdl/bookmark-sync/tree/main#command-line-interface)
- [Syncing Process](https://github.com/jonasrdl/bookmark-sync/tree/main#syncing-process)
- [Contributing](https://github.com/jonasrdl/bookmark-sync/tree/main#contributing)

### Installation
Bookmark Sync requires Go to be installed on your system. You can download and install Go from the official website: https://golang.org/dl/

To install Bookmark Sync, follow these steps:

1. Clone the repository:
`git clone https://github.com/jonasrdl/bookmark-sync.git`
2. Navigate to the project directory:
`cd bookmark-sync`
3. Build the executable:
`go build -o bookmark-sync cmd/bookmarksync/bookmarksync.go`
4. Add the executable to your system's PATH if desired.

### Usage
Bookmark Sync provides a command-line interface (CLI) that allows you to sync bookmarks between different browsers. Here's how you can use it:


`bookmark-sync sync --from=chromium --to=firefox`   
This command will sync bookmarks from Chromium to Firefox. You can replace chromium and firefox with your desired source and destination browsers.   
Currently available browsers are Chromium and Firefox. (More will follow).

### Command-Line Interface
Bookmark Sync CLI supports the following commands:

- `sync`: Sync bookmarks between browsers.
  - Flags:
    - --from: Source browser (chromium or firefox)
    - --to: Destination browser (chromium or firefox)
- `list-profiles`: List available profiles for a specific browser. **(Currently only for Firefox)**
  - Flags:
    - `--browser`: Browser for which to list profiles (e.g. firefox)
    
    
  
### Syncing Process
Bookmark Sync reads bookmarks from the source browser and merges them with the bookmarks in the destination browser. The merging process ensures that existing bookmarks are preserved, and new bookmarks are added.   
1. The source bookmarks are read from the source browser (Chromium or Firefox).
2. The destination bookmarks are read from the destination browser.
3. The bookmarks are merged, and duplicates are eliminated.
4. The merged bookmarks are written back to the destination browser.

### Listing Available Profiles
You can use the `list-profiles` command to list available profiles for a specific browser.
Please note that listing Chromium profiles is **currently** not supported in the current version due to limitations.
Chromium will follow in the future.

`bookmark-sync list-profiles --browser=firefox`   
This command will list all available profiles for Firefox browser.

### Contributing
Contributions are welcome! If you'd like to contribute to Bookmark Sync, please follow these steps:

1. Fork the repository on GitHub.
2. Create a new branch with a descriptive name.
3. Make your changes and test thoroughly.
4. Commit your changes with clear commit messages.
5. Push your branch to your forked repository.
6. Create a pull request to the main repository.
