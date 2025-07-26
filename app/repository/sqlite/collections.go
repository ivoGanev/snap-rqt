package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"snap-rq/app/entity"
	"snap-rq/app/repository"
	"time"
)

type SQLiteCollectionRepository struct {
	db *sql.DB
}


func NewCollectionRepository(sqliteDb *SQLiteDb) *SQLiteCollectionRepository {
	repo := &SQLiteCollectionRepository{
		db: sqliteDb.db,
	}

	// Ensure table exists
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS collections (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		created_at DATETIME NOT NULL,
		modified_at DATETIME,
		row_position INTEGER
	);`

	if _, err := repo.db.Exec(createTableQuery); err != nil {
		panic(fmt.Errorf("failed to create collections table: %w", err))
	}

	// Check if any collection exists
	var count int
	err := repo.db.QueryRow(`SELECT COUNT(*) FROM collections`).Scan(&count)
	if err != nil {
		panic(fmt.Errorf("failed to check collections count: %w", err))
	}

	// Insert default collection if none exists
	if count == 0 {
		defaultCol := entity.NewCollection("Default Collection", "Automatically created", 0)
		_, err := repo.db.Exec(`INSERT INTO collections (id, name, description, created_at, row_position) VALUES (?, ?, ?, ?, ?)`,
			defaultCol.Id, defaultCol.Name, defaultCol.Description, defaultCol.CreatedAt, defaultCol.RowPosition)
		if err != nil {
			panic(fmt.Errorf("failed to insert default collection: %w", err))
		}
	}

	return repo
}

func (m *SQLiteCollectionRepository) ShiftCollections(position int, direction string) error  {
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var query string
	switch direction {
	case repository.SHIFT_UP:
		query = `UPDATE collections SET row_position = row_position + 1 WHERE row_position >= ?`
	case repository.SHIFT_DOWN:
		query = `UPDATE collections SET row_position = row_position - 1 WHERE row_position >= ?`
	}

	_, err = tx.Exec(query, position)
	if err != nil {
		return fmt.Errorf("failed to shift requests: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit shift transaction: %w", err)
	}

	return nil
}

func (s *SQLiteCollectionRepository) GetCollections() ([]entity.Collection, error) {
	query := `SELECT id, name, description, created_at, modified_at FROM collections ORDER BY created_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query collections: %w", err)
	}
	defer rows.Close()

	var collections []entity.Collection
	for rows.Next() {
		var col entity.Collection
		var modifiedAt sql.NullTime
		err := rows.Scan(&col.Id, &col.Name, &col.Description, &col.CreatedAt, &modifiedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan collection: %w", err)
		}
		if modifiedAt.Valid {
			col.ModifiedAt = &modifiedAt.Time
		} else {
			col.ModifiedAt = nil
		}
		collections = append(collections, col)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return collections, nil
}

func (s *SQLiteCollectionRepository) CreateCollection(c *entity.Collection) error {
	query := `INSERT INTO collections (id, name, description, created_at, row_position) VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, c.Id, c.Name, c.Description, c.CreatedAt, c.RowPosition)
	if err != nil {
		return fmt.Errorf("failed to insert collection: %w", err)
	}
	return nil
}

func (s *SQLiteCollectionRepository) DeleteCollection(id string) error {
	res, err := s.db.Exec(`DELETE FROM collections WHERE id = ?`, id)

	if err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("collection not found")
	}
	return nil
}

func (s *SQLiteCollectionRepository) GetCollection(id string) (entity.Collection, error) {
	var col entity.Collection
	var modifiedAt sql.NullTime

	query := `SELECT id, name, description, created_at, modified_at, row_position FROM collections WHERE id = ?`

	err := s.db.QueryRow(query, id).Scan(&col.Id, &col.Name, &col.Description, &col.CreatedAt, &modifiedAt, &col.RowPosition)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Collection{}, errors.New("collection not found")
		}
		return entity.Collection{}, fmt.Errorf("failed to query collection: %w", err)
	}
	if modifiedAt.Valid {
		col.ModifiedAt = &modifiedAt.Time
	}
	return col, nil
}

func (s *SQLiteCollectionRepository) UpdateCollection(updated entity.Collection) (entity.Collection, error) {
	now := time.Now()
	updated.ModifiedAt = &now

	query := `
	UPDATE collections SET name=?, description=?, modified_at=?
	WHERE id=?
`

	res, err := s.db.Exec(query, updated.Name, updated.Description, updated.ModifiedAt, updated.Id)
	if err != nil {
		return entity.Collection{}, fmt.Errorf("failed to update collection: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return entity.Collection{}, fmt.Errorf("failed to get affected rows: %w", err)
	}
	if affected == 0 {
		return entity.Collection{}, errors.New("collection not found")
	}

	return updated, nil
}
