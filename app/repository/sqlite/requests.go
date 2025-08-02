package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"snap-rq/app/entity"
	logger "snap-rq/app/log"
	"snap-rq/app/repository"
)

const SQLITE_REQUEST_REPO_LOG_TAG = "[SQLite Requests Repository]"

type SQLiteRequestsRepository struct {
	db *sql.DB
}

func NewRequestsRepository(sqliteDb *SQLiteDb) *SQLiteRequestsRepository {
	repo := &SQLiteRequestsRepository{
		db: sqliteDb.db,
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS requests (
		id TEXT PRIMARY KEY,
		collection_id TEXT NOT NULL,
		name TEXT,
		description TEXT,
		method TEXT,
		url TEXT,
		headers TEXT, -- JSON
		body TEXT,
		created_at DATETIME,
		modified_at DATETIME,
		row_position INTEGER,
		FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
	);`

	if _, err := repo.db.Exec(createTableQuery); err != nil {
		panic(fmt.Errorf("failed to create requests table: %w", err))
	}

	return repo
}

// ShiftRequests updates the RowPosition of requests in a collection starting from a position, moving up or down
func (m *SQLiteRequestsRepository) ShiftRequests(collectionId string, startingPosition int, direction string) error {
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var query string
	switch direction {
	case repository.SHIFT_UP:
		query = `UPDATE requests SET row_position = row_position + 1 WHERE collection_id = ? AND row_position >= ?`
	case repository.SHIFT_DOWN:
		query = `UPDATE requests SET row_position = row_position - 1 WHERE collection_id = ? AND row_position >= ?`
	}

	_, err = tx.Exec(query, collectionId, startingPosition)
	if err != nil {
		return fmt.Errorf("failed to shift requests: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit shift transaction: %w", err)
	}

	return nil
}

func (m *SQLiteRequestsRepository) DeleteRequest(id string) error {
	res, err := m.db.Exec(`DELETE FROM requests WHERE id = ?`, id)

	if err != nil {
		return fmt.Errorf("failed to delete request: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("request id not found: %s", id)
	}
	return nil
}

func (m *SQLiteRequestsRepository) CreateRequest(request entity.Request) error {
	query := `
	INSERT INTO requests (id, name, description, modified_at, created_at, collection_id, row_position, method, url, headers, body)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	_, err := m.db.Exec(query,
		request.Id,
		request.Name,
		request.Description,
		request.ModifiedAt,
		request.CreatedAt,
		request.CollectionID,
		request.RowPosition,
		request.Method,
		request.Url,
		request.Headers,
		request.Body,
	)
	if err != nil {
		logger.Println("failed to create request:", err)
	}
	return nil
}

// GetRequestsBasic returns minimal request info for a collection
func (m *SQLiteRequestsRepository) GetRequestsBasic(collectionId string) ([]entity.RequestBasic, error) {
	query := `
	SELECT id, name, row_position, method, url
	FROM requests
	WHERE collection_id = ?
	ORDER BY row_position ASC
`

	rows, err := m.db.Query(query, collectionId)
	if err != nil {
		return nil, fmt.Errorf("failed to query basic requests: %w", err)
	}
	defer rows.Close()

	var items []entity.RequestBasic
	for rows.Next() {
		var rb entity.RequestBasic
		err := rows.Scan(
			&rb.Id,
			&rb.Name,
			&rb.RowPosition,
			&rb.Method,
			&rb.Url,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan request basic: %w", err)
		}
		items = append(items, rb)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return items, nil
}

func (m *SQLiteRequestsRepository) GetRequest(id string) (entity.Request, error) {
	var r entity.Request

	query := `
	SELECT id, collection_id, name, description, method, url, headers, body, created_at, modified_at, row_position
	FROM requests
	WHERE id = ?
	`

	err := m.db.QueryRow(query, id).Scan(
		&r.Id,
		&r.CollectionID,
		&r.Name,
		&r.Description,
		&r.Method,
		&r.Url,
		&r.Headers,
		&r.Body,
		&r.CreatedAt,
		&r.ModifiedAt,
		&r.RowPosition,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Request{}, fmt.Errorf("request id not found: %s", id)
		}
		return entity.Request{}, fmt.Errorf("failed to get request: %w", err)
	}

	return r, nil
}

func (m *SQLiteRequestsRepository) UpdateRequest(updated entity.Request) (entity.Request, error) {
	query := `
	UPDATE requests
	SET collection_id = ?, name = ?, row_position = ?, method = ?, url = ?, headers = ?, body = ?
	WHERE id = ?
`

	res, err := m.db.Exec(query,
		updated.CollectionID,
		updated.Name,
		updated.RowPosition,
		updated.Method,
		updated.Url,
		updated.Headers,
		updated.Body,
		updated.Id,
	)
	logger.Info(SQLITE_REQUEST_REPO_LOG_TAG, "Trying to update request with Id:", updated.Id, "The updated data is ", updated)

	if err != nil {
		return entity.Request{}, fmt.Errorf("failed to update request: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return entity.Request{}, fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return entity.Request{}, fmt.Errorf("request id not found: %s", updated.Id)
	}

	return updated, nil
}
