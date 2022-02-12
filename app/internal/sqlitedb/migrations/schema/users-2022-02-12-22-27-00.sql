CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT CHECK ( role IN('admin', 'reader') ) NOT NULL,
    token TEXT DEFAULT '' NOT NULL
);