package sql_migrations

import (
	"context"
	"database/sql"
	"errors"
)

var migrations = []string{
	`CREATE TABLE IF NOT EXISTS beacon_indicies (
		beacon_block_root BLOB NOT NULL CHECK(length(beacon_block_root) = 32),
		slot INTEGER NOT NULL,
		state_root BLOB NOT NULL CHECK(length(state_root) = 32),
		parent_block_root BLOB NOT NULL CHECK(length(parent_block_root) = 32),
		canonical INTEGER NOT NULL DEFAULT 0, -- 0 for false, 1 for true
		PRIMARY KEY (beacon_block_root)
	);`,
	`CREATE INDEX idx_slot ON beacon_indicies (slot);`,
	`CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_canonical 
		ON beacon_indicies (slot) 
		WHERE canonical = 1;`,
	`CREATE TABLE IF NOT EXISTS data_config (prune_depth INTEGER NOT NULL)`,
}

func ApplyMigrations(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS schema_version (version INTEGER NOT NULL);"); err != nil {
		return err
	}
	// Get current schema version
	var currentVersion int
	err := tx.QueryRow("SELECT version FROM schema_version").Scan(&currentVersion)
	if errors.Is(err, sql.ErrNoRows) {
		currentVersion = -1
	} else if err != nil {
		return err
	}

	// Apply missing migrations
	for i := currentVersion + 1; i < len(migrations); i++ {
		_, err = tx.Exec(migrations[i])
		if err != nil {
			return err
		}

		// Update schema version
		_, err = tx.Exec("UPDATE schema_version SET version = ?", i)
		if err != nil {
			return err
		}
	}

	return nil
}
