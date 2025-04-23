package authz

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func ApplyMigrations(db *sql.DB, migrationsDir string) error {
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS authz_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, file := range files {
		version := filepath.Base(file)
		applied, err := isMigrationApplied(db, version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		fmt.Println("📥 Applying migration:", version)

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to apply %s: %w", version, err)
		}

		if _, err := db.ExecContext(ctx,
			`INSERT INTO authz_migrations (version) VALUES ($1)`, version,
		); err != nil {
			return err
		}
	}

	fmt.Println("✅ Migrations applied successfully.")
	return nil
}

func isMigrationApplied(db *sql.DB, version string) (bool, error) {
	var exists bool
	err := db.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM authz_migrations WHERE version = $1
	)`, version).Scan(&exists)
	return exists, err
}
