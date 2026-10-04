package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateCategoriesTable, downCreateCategoriesTable)
}

func upCreateCategoriesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
        CREATE TABLE categories (
            id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
            slug TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL UNIQUE
        )
    `)

	return err
}

func downCreateCategoriesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
        DROP TABLE IF EXISTS categories
    `)

	return err
}
