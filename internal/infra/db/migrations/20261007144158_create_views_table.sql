-- +goose Up
CREATE TABLE IF NOT EXISTS views (
    article_id BIGINT NOT NULL,
    views BIGINT NOT NULL,

    PRIMARY KEY (article_id),

    FOREIGN KEY (article_id) REFERENCES articles(id)
);

-- +goose Down
DROP TABLE IF EXISTS views;
