-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS todos (
    id CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    completed tinyint unsigned NOT NULL DEFAULT 0,
    created_at datetime NOT NULL,
    updated_at datetime NOT NULL,
    PRIMARY KEY (`id`)
) ROW_FORMAT=DYNAMIC;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS todos;
