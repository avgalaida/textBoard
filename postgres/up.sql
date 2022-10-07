DROP TABLE IF EXISTS posts;
CREATE TABLE posts (
                       id VARCHAR(32) PRIMARY KEY,
                       body TEXT NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE NOT NULL
);