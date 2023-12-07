CREATE TABLE IF NOT EXISTS login (
	id uuid PRIMARY KEY REFERENCES credential_record ON UPDATE CASCADE,
	username text,
	password text,
	url text
);
