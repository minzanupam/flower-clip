CREATE TABLE users (
	id INTEGER PRIMARY KEY autoincrement,
	fullname TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL
);
