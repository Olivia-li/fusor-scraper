CREATE TABLE threads (
    id SERIAL PRIMARY KEY,
    post_id VARCHAR(255) UNIQUE,
    thread_id VARCHAR(255),
    title TEXT,
    content TEXT,
    author VARCHAR(255),
    post_time TIMESTAMP
);