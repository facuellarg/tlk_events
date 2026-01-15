package repository

import (
	"challenge/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// Database holds the database connection
type Database struct {
	DB *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(ctx context.Context, dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	// Optional: Configure connection pool settings
	db.SetMaxOpenConns(1) // SQLite works best with single connection for writes
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	// Enable foreign keys and WAL mode for better concurrency
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	_, err = db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	log.Println("Successfully connected to SQLite database")

	return &Database{DB: db}, nil
}

// Close closes the database connection
func (db *Database) Close() {
	db.DB.Close()
	log.Println("Database connection closed")
}

// CreateTable creates the events table if it doesn't exist
func (db *Database) CreateTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL CHECK(length(title) <= 100),
		description TEXT,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_events_start_time ON events(start_time);
	CREATE INDEX IF NOT EXISTS idx_events_end_time ON events(end_time);
	`

	_, err := db.DB.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	log.Println("Table 'events' is ready")
	return nil
}

// InsertEvent inserts a new event into the database
func (db *Database) InsertEvent(ctx context.Context, event *models.Event) error {
	// Generate UUID if not provided
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}

	// Set created_at if not provided
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO events (id, title, description, start_time, end_time, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.DB.ExecContext(ctx, query,
		event.ID.String(),
		event.Title,
		event.Description,
		event.StartTime.Format(time.RFC3339),
		event.EndTime.Format(time.RFC3339),
		event.CreatedAt.Format(time.RFC3339),
	)

	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}

	log.Printf("Event inserted successfully with ID: %s", event.ID)
	return nil
}

// GetEventByID retrieves an event by its ID
func (db *Database) GetEventByID(ctx context.Context, id uuid.UUID) (*models.Event, error) {
	query := `
		SELECT id, title, description, start_time, end_time, created_at
		FROM events
		WHERE id = ?
	`

	var event models.Event
	var idStr string
	var startTimeStr, endTimeStr, createdAtStr string

	err := db.DB.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr,
		&event.Title,
		&event.Description,
		&startTimeStr,
		&endTimeStr,
		&createdAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Parse UUID
	event.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse UUID: %w", err)
	}

	// Parse timestamps
	event.StartTime, err = time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start_time: %w", err)
	}

	event.EndTime, err = time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end_time: %w", err)
	}

	event.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	return &event, nil
}

// GetAllEvents retrieves all events from the database
func (db *Database) GetAllEvents(ctx context.Context) ([]*models.Event, error) {
	query := `
		SELECT id, title, description, start_time, end_time, created_at
		FROM events
		ORDER BY start_time ASC
	`

	rows, err := db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []*models.Event
	for rows.Next() {
		var event models.Event
		var idStr string
		var startTimeStr, endTimeStr, createdAtStr string

		err := rows.Scan(
			&idStr,
			&event.Title,
			&event.Description,
			&startTimeStr,
			&endTimeStr,
			&createdAtStr,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		// Parse UUID
		event.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse UUID: %w", err)
		}

		// Parse timestamps
		event.StartTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse start_time: %w", err)
		}

		event.EndTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse end_time: %w", err)
		}

		event.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}

		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return events, nil
}

// UpdateEvent updates an existing event
func (db *Database) UpdateEvent(ctx context.Context, event *models.Event) error {
	query := `
		UPDATE events
		SET title = ?, description = ?, start_time = ?, end_time = ?
		WHERE id = ?
	`

	result, err := db.DB.ExecContext(ctx, query,
		event.Title,
		event.Description,
		event.StartTime.Format(time.RFC3339),
		event.EndTime.Format(time.RFC3339),
		event.ID.String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event not found")
	}

	log.Printf("Event updated successfully with ID: %s", event.ID)
	return nil
}

// DeleteEvent deletes an event by ID
func (db *Database) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = ?`

	result, err := db.DB.ExecContext(ctx, query, id.String())
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event not found")
	}

	log.Printf("Event deleted successfully with ID: %s", id)
	return nil
}

// Example usage
func main() {
	ctx := context.Background()

	// Get database path from environment variable
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./events.db"
	}

	// Create database connection
	db, err := NewDatabase(ctx, dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create table
	if err := db.CreateTable(ctx); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Example: Insert a new event
	description := "Team meeting to discuss Q1 goals"
	newEvent := &models.Event{
		Title:       "Team Meeting",
		Description: &description,
		StartTime:   time.Now().Add(24 * time.Hour),
		EndTime:     time.Now().Add(25 * time.Hour),
	}

	if err := db.InsertEvent(ctx, newEvent); err != nil {
		log.Printf("Failed to insert event: %v", err)
	}

	// Example: Get all events
	events, err := db.GetAllEvents(ctx)
	if err != nil {
		log.Printf("Failed to get events: %v", err)
	} else {
		log.Printf("Found %d events", len(events))
		for _, event := range events {
			log.Printf("Event: %s - %s", event.ID, event.Title)
		}
	}
}
