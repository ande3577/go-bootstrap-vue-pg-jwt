-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login varchar(255) UNIQUE NOT NULL,
    email varchar(255) UNIQUE NOT NULL,
    hashed_password varchar(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;