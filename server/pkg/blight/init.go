package blight

import (
	"database/sql"
	"fmt"
)

func (c *Client) initializeDatabase() error {
	if err := c.createBlobsTable(); err != nil {
		return fmt.Errorf("%w: %v", ErrCreatingDatabase, err)
	}
	if err := c.createUpdateTrigger(); err != nil {
		return fmt.Errorf("%w: %v", ErrCreatingDatabase, err)
	}
	return nil
}

func (c *Client) createBlobsTable() error {
	stmt := `
		CREATE TABLE IF NOT EXISTS blobs (
			path TEXT PRIMARY KEY,
			blob BLOB NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`

	_, err := c.DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("failed to create blobs table: %w", err)
	}
	return nil
}

func (c *Client) createUpdateTrigger() error {
	stmt := `
		CREATE TRIGGER IF NOT EXISTS tr_blobs_set_updated_at
		AFTER UPDATE ON blobs
		FOR EACH ROW
		BEGIN
			UPDATE blobs
			SET updated_at = CURRENT_TIMESTAMP
			WHERE path = NEW.path;
		END;`

	_, err := c.DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("failed to create update trigger: %w", err)
	}
	return nil
}

func (c *Client) initializePreparedStatements() error {
	add, err := c.DB.Prepare(addStmt)
	if err != nil {
		return err
	}
	get, err := c.DB.Prepare(getStmt)
	if err != nil {
		return err
	}
	del, err := c.DB.Prepare(deleteStmt)
	if err != nil {
		return err
	}

	c.statements = struct {
		add    *sql.Stmt
		get    *sql.Stmt
		delete *sql.Stmt
	}{
		add:    add,
		get:    get,
		delete: del,
	}
	return nil
}
