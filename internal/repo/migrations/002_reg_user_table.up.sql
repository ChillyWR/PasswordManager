CREATE TABLE IF NOT EXISTS reg_user (
	id uuid PRIMARY KEY,
	name text NOT NULL,
	password text NOT NULL,
	created_on timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_on timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
