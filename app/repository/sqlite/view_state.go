package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"snap-rq/app/entity"
	logger "snap-rq/app/log"
)

type SQLiteViewStateRepository struct {
	db *sql.DB
}

const VIEW_STATE_LOG_TAG = "[SQLite State Repository]"

func NewStateRepository(sqliteDb *SQLiteDb, c *SQLiteCollectionRepository, r *SQLiteRequestsRepository) *SQLiteViewStateRepository {
	repo := &SQLiteViewStateRepository{db: sqliteDb.db}

	// Ensure table exists
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS app_view_state (
		id INTEGER PRIMARY KEY,
		focused_collection_id TEXT,
		focused_view TEXT,
		focused_request_ids TEXT -- stored as JSON
	);`
	if _, err := repo.db.Exec(createTableQuery); err != nil {
		panic(fmt.Errorf("failed to create app_view_state table: %w", err))
	}

	// Check if state exists
	var count int
	err := sqliteDb.db.QueryRow(`SELECT COUNT(*) FROM app_view_state WHERE id = 1`).Scan(&count)
	if err != nil {
		panic(fmt.Errorf("failed to query app_view_state: %w", err))
	}

	if count == 0 {
		// No saved state yet — create default
		collections, err := c.GetCollections()
		if err != nil || len(collections) == 0 {
			panic("failed to initialize default view state")
		}

		selectedRequests := make(map[string]string)
		for _, collection := range collections {
			req, _ := r.GetRequestsBasic(collection.Id)
			if len(req) > 0 {
				selectedRequests[collection.Id] = req[0].Id
			}
		}

		state := entity.AppViewState{
			FocusedRequestIds:   selectedRequests,
			FocusedCollectionId: collections[0].Id,
			FocusedView:         "requests",
		}

		if err := repo.SetState(state); err != nil {
			panic(fmt.Errorf("failed to save initial app_view_state: %w", err))
		}
	}

	return repo
}

func (m *SQLiteViewStateRepository) SetState(state entity.AppViewState) error {
	reqIdsJSON, err := json.Marshal(state.FocusedRequestIds)
	if err != nil {
		return fmt.Errorf("failed to encode request ids: %w", err)
	}

	_, err = m.db.Exec(`
	INSERT INTO app_view_state (id, focused_view, focused_collection_id, focused_request_ids)
	VALUES (1, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		focused_view = excluded.focused_view,
		focused_collection_id = excluded.focused_collection_id,
		focused_request_ids = excluded.focused_request_ids
	`, state.FocusedView, state.FocusedCollectionId, string(reqIdsJSON))

	if err != nil {
		return fmt.Errorf("failed to store view state: %w", err)
	}

	logger.Println(VIEW_STATE_LOG_TAG, "Set state:", state)
	return nil
}

func (m *SQLiteViewStateRepository) GetState() (entity.AppViewState, error) {
	var view, collectionId, requestIdsJSON string

	err := m.db.QueryRow(`
		SELECT focused_view, focused_collection_id, focused_request_ids
		FROM app_view_state
		WHERE id = 1
	`).Scan(&view, &collectionId, &requestIdsJSON)

	if err != nil {
		if err == sql.ErrNoRows {
			return entity.AppViewState{}, nil // Or return error if required
		}
		return entity.AppViewState{}, err
	}

	var requestIds map[string]string
	if err := json.Unmarshal([]byte(requestIdsJSON), &requestIds); err != nil {
		return entity.AppViewState{}, fmt.Errorf("failed to decode request ids: %w", err)
	}

	state := entity.AppViewState{
		FocusedView:         view,
		FocusedCollectionId: collectionId,
		FocusedRequestIds:   requestIds,
	}

	return state, nil
}
