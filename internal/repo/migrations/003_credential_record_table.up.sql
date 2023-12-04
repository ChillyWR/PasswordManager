CREATE TABLE IF NOT EXISTS credential_record (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	name varchar NOT NULL,
	notes text,
	created_on timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_on timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by uuid NOT NULL REFERENCES reg_user ON UPDATE CASCADE,
	updated_by uuid NOT NULL REFERENCES reg_user ON UPDATE CASCADE
);
