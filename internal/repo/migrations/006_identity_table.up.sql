CREATE TABLE IF NOT EXISTS identity (
	id uuid PRIMARY KEY REFERENCES credential_record ON UPDATE CASCADE,
	first_name varchar,
	middle_name varchar,
	last_name varchar,
	address varchar,
	email varchar,
	phone_number varchar,
	passport_number varchar,
	country varchar
);
