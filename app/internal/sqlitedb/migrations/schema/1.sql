CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT CHECK ( role IN('admin', 'reader') ) NOT NULL,
    token TEXT DEFAULT '' NOT NULL
);

CREATE TABLE IF NOT EXISTS onecdbs (
    id TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    server TEXT NOT NULL,
    ref TEXT NOT NULL,
    description TEXT NOT NULL,
    folder TEXT NOT NULL
);
