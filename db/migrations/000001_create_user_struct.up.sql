CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    guid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ip VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    refresh_token VARCHAR(255),
    access_token UUID DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS used_tokens (
    id SERIAL PRIMARY KEY,
    user_guid VARCHAR(36) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    used_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);