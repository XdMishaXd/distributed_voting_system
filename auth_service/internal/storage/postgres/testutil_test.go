package postgres

import (
	"database/sql"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pressly/goose/v3"
)

const targetMigrationVersion = 20251119160405

func migrationsPath(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to determine caller for migrations path")
	}
	root := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..")
	return filepath.Join(root, "migrations", "postgres")
}

// applyMigrations поднимает схему перед тестом и откатывает после — чтобы
// один race-тест не оставлял состояние, которое сломает следующий.
func applyMigrations(t *testing.T, db *sql.DB) {
	t.Helper()

	dir := migrationsPath(t)
	if err := goose.UpTo(db, dir, targetMigrationVersion); err != nil {
		t.Fatalf("goose.Up (dir=%s): %v", dir, err)
	}

	t.Cleanup(func() {
		if err := goose.DownTo(db, dir, 0); err != nil {
			t.Logf("goose.Down cleanup failed (dir=%s): %v", dir, err)
		}
	})
}
