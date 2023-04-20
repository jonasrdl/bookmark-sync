const PORT = 6758;
const SERVER_URL = `http://localhost:${PORT}`;
const browser = chrome || browser;
const SYNCED_FOLDER_NAME = 'Synced Bookmarks';

async function syncBookmark(action, bookmark) {
  try {
    const response = await fetch(`${SERVER_URL}/sync`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action, bookmark }),
    });

    if (!response.ok) {
      throw new Error(`Server error: ${response.statusText}`);
    }
  } catch (error) {
    console.error(`Failed to sync bookmark: ${error.message}`);
  }
}

async function copyBookmarks(srcBookmarks, parentId) {
  for (const bookmark of srcBookmarks) {
    if (bookmark.url) {
      await browser.bookmarks.create({ parentId, title: bookmark.title, url: bookmark.url });
    } else {
      const newFolder = await browser.bookmarks.create({ parentId, title: bookmark.title });
      await copyBookmarks(bookmark.children, newFolder.id);
    }
  }
}

async function initialUpload(bookmarks) {
  try {
    const response = await fetch(`${SERVER_URL}/initial-upload`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ bookmarks }),
    });

    if (!response.ok) {
      throw new Error(`Server error: ${response.statusText}`);
    }

    // Create a new bookmark folder for synced bookmarks
    const syncedFolder = await browser.bookmarks.create({
      title: SYNCED_FOLDER_NAME,
    });

    // Copy bookmarks to the new folder
    await copyBookmarks(bookmarks[0].children, syncedFolder.id);
  } catch (error) {
    console.error(`Failed to upload initial bookmarks: ${error.message}`);
  }
}

// Add event listeners for bookmark changes
browser.bookmarks.onCreated.addListener((id, bookmark) => {
  syncBookmark('created', bookmark);
});

browser.bookmarks.onRemoved.addListener((id, removeInfo) => {
  syncBookmark('removed', { id });
});

browser.bookmarks.onChanged.addListener((id, changeInfo) => {
  syncBookmark('changed', { id, ...changeInfo });
});

// Perform initial upload of all bookmarks
(async function () {
  const bookmarks = await browser.bookmarks.getTree();
  initialUpload(bookmarks);
})();
