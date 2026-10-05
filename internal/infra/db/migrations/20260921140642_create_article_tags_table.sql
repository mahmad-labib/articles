-- +goose Up
CREATE TABLE IF NOT EXISTS article_tags (
    article_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,

    PRIMARY KEY (article_id, tag_id),

    FOREIGN KEY (article_id) REFERENCES articles(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);

-- +goose Down
DROP TABLE IF EXISTS article_tags;
