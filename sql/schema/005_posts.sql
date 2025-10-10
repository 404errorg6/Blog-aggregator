-- +goose Up
CREATE TABLE posts (
  id UUID UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  title TEXT NOT NULL,
  url TEXT UNIQUE NOT NULL,
  description TEXT NOT NULL, 
  published_at TIMESTAMP NOT NULL,
  feed_id UUID NOT NULL
);

-- +goose Down
DROP TABLE posts;
