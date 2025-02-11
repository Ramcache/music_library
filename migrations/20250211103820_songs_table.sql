-- +goose Up
CREATE TABLE IF NOT EXISTS songs (
                                     id SERIAL PRIMARY KEY,
                                     group_name VARCHAR(255),
    song VARCHAR(255),
    release_date VARCHAR(50),
    text TEXT,
    link TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- +goose Down
DROP TABLE songs;
