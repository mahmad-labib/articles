package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateArticleTagsTable, downCreateArticleTagsTable)
}

func upCreateArticleTagsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
        CREATE TABLE article_tags (
			article_id BIGINT NOT NULL,
			tag_id BIGINT NOT NULL,

			PRIMARY KEY (article_id, tag_id),

			FOREIGN KEY (article_id) REFERENCES articles(id),
			FOREIGN KEY (tag_id) REFERENCES tags(id)
        )
    `)
	return err
}

func downCreateArticleTagsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS article_tags
	`)
	return err
}
