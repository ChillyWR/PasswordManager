CREATE TABLE IF NOT EXISTS reg_user (
	id uuid PRIMARY KEY,
	username varchar NOT NULL,
	password varchar NOT NULL,
	email varchar
);
