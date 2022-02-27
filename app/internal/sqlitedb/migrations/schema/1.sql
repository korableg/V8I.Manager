CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT CHECK ( role IN('admin', 'reader') ) NOT NULL,
    token TEXT DEFAULT '' NOT NULL
);

CREATE TABLE IF NOT EXISTS onecdbs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT UNIQUE NOT NULL,
    name TEXT DEFAULT '' NOT NULL,
    connect TEXT NOT NULL,
    order_in_list INTEGER DEFAULT 0 NOT NULL,
    order_in_tree INTEGER DEFAULT 0 NOT NULL,
    folder TEXT DEFAULT '' NOT NULL,
    client_connection_speed TEXT CHECK ( client_connection_speed IN ('normal', 'low') ) DEFAULT 'normal' NOT NULL,
    app TEXT CHECK (app IN ('Auto', 'ThinClient', 'ThickClient', 'WebClient')) DEFAULT 'Auto' NOT NULL,
    wa INTEGER CHECK (wa IN (1, 0)) DEFAULT 1 NOT NULL,
    version TEXT DEFAULT '' NOT NULL,
    additional_parameters TEXT DEFAULT '' NOT NULL
);
