package db

import (
	"embed"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Migration represents a database migration manager that uses goose and gorm.DB.
type Migration struct {
	DB *gorm.DB
}

// NewMigration creates a new Migration instance with the given gorm.DB connection.
func NewMigration(db *gorm.DB) (*Migration, error) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}

	return &Migration{
		DB: db,
	}, nil
}

// Up applies all available migrations.
func (m *Migration) Up() error {
	sqlDB, err := m.DB.DB()
	if err != nil {
		return err
	}

	return goose.Up(sqlDB, "migrations")
}

// Down rolls back the most recent migration.
func (m *Migration) Down() error {
	sqlDB, err := m.DB.DB()
	if err != nil {
		return err
	}

	return goose.Down(sqlDB, "migrations")
}

// Reset rolls back all migrations and then applies them again.
func (m *Migration) Reset() error {
	sqlDB, err := m.DB.DB()
	if err != nil {
		return err
	}

	return goose.Reset(sqlDB, "migrations")
}
