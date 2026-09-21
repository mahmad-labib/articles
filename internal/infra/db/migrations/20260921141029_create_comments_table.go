package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateCommentsTable, downCreateCommentsTable)
}

func upCreateCommentsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
        CREATE TABLE comments (
            id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
            article_id BIGINT NOT NULL,
            author_id BIGINT NOT NULL,
            body TEXT NOT NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

            FOREIGN KEY (article_id) REFERENCES articles(id),
            FOREIGN KEY (author_id) REFERENCES users(id)
        )
    `)
	return err
}

func downCreateCommentsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS comments
	`)
	return err
}
