package storage

const schemaSQL = `
CREATE TABLE IF NOT EXISTS blog (
  id INTEGER PRIMARY KEY,
  url TEXT,
  title TEXT,
  html TEXT,
  tags TEXT
);
`
