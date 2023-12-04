CREATE TABLE IF NOT EXISTS card (
	id uuid PRIMARY KEY REFERENCES credential_record ON UPDATE CASCADE,
	brand varchar,
	number varchar,
	expiration_month varchar,
	expiration_year varchar,
	cvv varchar
);
