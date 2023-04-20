import { JSONFileSync } from 'lowdb/node';
import { Low } from 'lowdb';
import express from 'express';
import bodyParser from 'body-parser';
import cors from 'cors';

const app = express();
const PORT = process.env.PORT || 6758;

app.use(cors());

const adapter = new JSONFileSync('db.json');
const db = new Low(adapter, null);

db.read();
db.data = db.data || { bookmarks: [] };

app.use(bodyParser.json({ limit: '10mb' }));

app.post('/sync', (req, res) => {
  console.log("sync request", req);
  const { action, bookmark } = req.body;

  switch (action) {
    case 'created':
      db.data.bookmarks.push(bookmark);
      break;
    case 'removed':
      db.data.bookmarks = db.data.bookmarks.filter(b => b.id !== bookmark.id);
      break;
    case 'changed':
      const index = db.data.bookmarks.findIndex(b => b.id === bookmark.id);
      if (index !== -1) {
        db.data.bookmarks[index] = { ...db.data.bookmarks[index], ...bookmark };
      }
      break;
    default:
      res.status(400).json({ message: 'Invalid action' });
      return;
  }

  db.write();
  res.status(200).json({ message: 'Bookmark synced successfully' });
});

app.post('/initial-upload', (req, res) => {
  console.log("initial upload req", req);
  const { bookmarks } = req.body;
  db.data.bookmarks = bookmarks;
  db.write();
  res.status(200).json({ message: 'Initial bookmarks uploaded successfully' });
});

app.listen(PORT, () => {
  console.log(`Server running at http://localhost:${PORT}`);
});
