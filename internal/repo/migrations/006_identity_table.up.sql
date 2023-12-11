CREATE TABLE IF NOT EXISTS identity (
	id uuid PRIMARY KEY REFERENCES credential_record ON UPDATE CASCADE ON DELETE CASCADE,
	first_name text,
	middle_name text,
	last_name text,
	address text,
	email text,
	phone_number text,
	passport_number text,
	country text
);
