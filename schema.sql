CREATE TABLE users (
	id INTEGER PRIMARY KEY autoincrement,
	fullname TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	created_at TEXT NOT NULL
);

CREATE TABLE svgs (
	id INTEGER PRIMARY KEY autoincrement,
	name TEXT NOT NULL,
	file TEXT NOT NULL,
	created_at TEXT NOT NULL,
	user_id INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id)
);
