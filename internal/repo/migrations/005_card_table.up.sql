CREATE TABLE IF NOT EXISTS card (
	id uuid PRIMARY KEY REFERENCES credential_record ON UPDATE CASCADE ON DELETE CASCADE,
	brand text,
	number text,
	expiration_month text,
	expiration_year text,
	cvv text
);
