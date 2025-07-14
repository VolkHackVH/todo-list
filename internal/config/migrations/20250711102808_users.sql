-- migrate:up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- migrate:down
DROP TABLE IF EXISTS users;

