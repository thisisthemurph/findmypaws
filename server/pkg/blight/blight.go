package blight

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrCreatingDatabase = errors.New("error creating database")
	ErrBlobNotFound     = errors.New("blob not found")
)

const addStmt = `
	insert into blobs (path, blob) 
	values (?, ?)
	on conflict(path) do update
		set blob = ?;`

const getStmt = "select path, blob, created_at, updated_at from blobs where path = ?;"

const deleteStmt = "delete from blobs where path = ?;"

type Client struct {
	DB         *sql.DB
	statements struct {
		add    *sql.Stmt
		get    *sql.Stmt
		delete *sql.Stmt
	}
}

type BlobResult struct {
	Path      string
	BLOB      io.Reader
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(path string) (*Client, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	c := &Client{DB: db}
	if err := c.initializeDatabase(); err != nil {
		db.Close()
		return nil, err
	}

	if err := c.initializePreparedStatements(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error preparing statements: %w", err)
	}

	return c, nil
}

func (c *Client) Close() error {
	if c.statements.add != nil {
		c.statements.add.Close()
	}
	if c.statements.get != nil {
		c.statements.add.Close()
	}
	if c.statements.delete != nil {
		c.statements.add.Close()
	}
	return c.DB.Close()
}

func (c *Client) Add(path string, r io.Reader) error {
	blob, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read blob data: %w", err)
	}

	if _, err := c.statements.add.Exec(path, blob, blob); err != nil {
		return fmt.Errorf("failed to add blob: %w", err)
	}
	return nil
}

func (c *Client) Get(path string) (*BlobResult, error) {
	var res BlobResult
	var blob []byte
	if err := c.statements.get.QueryRow(path).Scan(&res.Path, &blob, &res.CreatedAt, &res.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBlobNotFound
		}
		return nil, fmt.Errorf("failed to get blob: %w", err)
	}

	res.BLOB = bytes.NewReader(blob)
	return &res, nil
}

func (c *Client) Delete(path string) error {
	result, err := c.statements.delete.Exec(path)
	if err != nil {
		return fmt.Errorf("failed to delete blob: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrBlobNotFound
	}
	return nil
}
