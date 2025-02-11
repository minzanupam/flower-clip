CREATE TABLE users (
	id INTEGER PRIMARY KEY autoincrement,
	fullname TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL
);

CREATE TABLE svgs (
	id INTEGER PRIMARY KEY autoincrement,
	name TEXT NOT NULL,
	file BLOB NOT NULL,
	created_at TEXT NOT NULL
	-- add user identifier later
);
