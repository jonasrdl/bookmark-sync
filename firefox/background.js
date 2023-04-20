const browser = typeof chrome !== 'undefined' ? chrome : browser;

const SERVER_URL = 'http://localhost:6758';
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
  console.log('Initial upload called with:', bookmarks);

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
    await copyBookmarks(bookmarks, syncedFolder.id);
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
browser.bookmarks.getTree().then((bookmarks) => {
  const ids = [];
  bookmarks.forEach((bookmark) => {
    ids.push(bookmark.id);
    if (bookmark.children) {
      bookmark.children.forEach((childBookmark) => {
        ids.push(childBookmark.id);
      });
    }
  });
  browser.bookmarks.get(ids).then((bookmarks) => {
    console.log('Bookmarks:', bookmarks);
    initialUpload(bookmarks);
  });
}).catch((error) => {
  console.error(`Failed to get bookmarks: ${error.message}`);
});