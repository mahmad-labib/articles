package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateArticlesTable, downCreateArticlesTable)
}

func upCreateArticlesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
        CREATE TABLE articles (
            id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
            title TEXT NOT NULL,
            slug TEXT NOT NULL UNIQUE,
            body TEXT NOT NULL,
            status TEXT NOT NULL DEFAULT 'draft'
                CHECK (status IN ('draft', 'scheduled', 'published')),
            user_id BIGINT NOT NULL,
            category_id BIGINT NOT NULL,
            published_at TIMESTAMPTZ NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
            deleted_at TIMESTAMPTZ NULL,

            CONSTRAINT fk_articles_author
                FOREIGN KEY (user_id) REFERENCES users(id),

            CONSTRAINT fk_articles_category
                FOREIGN KEY (category_id) REFERENCES categories(id)
        )
    `)
	return err
}

func downCreateArticlesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
        DROP TABLE IF EXISTS articles
    `)
	return err
}
