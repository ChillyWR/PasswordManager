CREATE TABLE IF NOT EXISTS users
(
    id uuid primary key,
    name text not null,
    email text not null,
    login text not null,
    password text not null,
    phone text,
    created_at timestamp,
    updated_at timestamp
);
