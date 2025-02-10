CREATE TABLE users (
	id integer not null,
	fullname text not null,
	email text not null unique,
	password text not null
);
