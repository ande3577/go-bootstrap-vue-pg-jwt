
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE sessions (
	id SERIAL PRIMARY KEY,
	user_id int,
	FOREIGN KEY (user_id) REFERENCES users(id),
	session varchar(255) UNIQUE NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE sessions;
