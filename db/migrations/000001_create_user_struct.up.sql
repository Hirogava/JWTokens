CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    guid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ip VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    refresh_token VARCHAR(255),
    access_token UUID DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);