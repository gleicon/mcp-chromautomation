package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gleicon/mcp-chromautomation/internal/browser"
	_ "modernc.org/sqlite"
)

// Manager handles local data storage
type Manager struct {
	db   *sql.DB
	path string
}

// Session represents a saved browser session
type Session struct {
	ID        int                    `json:"id"`
	Name      string                 `json:"name"`
	Data      *browser.SessionData   `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// RequestLog represents a logged HTTP request
type RequestLog struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`
	Method    string    `json:"method"`
	Headers   string    `json:"headers"`
	Body      string    `json:"body"`
	Response  string    `json:"response"`
	Status    int       `json:"status"`
	Duration  int64     `json:"duration_ms"`
	Timestamp time.Time `json:"timestamp"`
}

// Screenshot represents a saved screenshot
type Screenshot struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	FilePath    string    `json:"file_path"`
	Size        int64     `json:"size"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// New creates a new storage manager
func New() *Manager {
	homeDir, _ := os.UserHomeDir()
	dbPath := filepath.Join(homeDir, ".mcp-chromautomation", "data.db")
	
	return &Manager{
		path: dbPath,
	}
}

// Init initializes the storage manager
func (m *Manager) Init() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(m.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	
	// Open database
	db, err := sql.Open("sqlite", m.path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	m.db = db
	
	// Create tables
	if err := m.createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}
	
	return nil
}

// Close closes the storage manager
func (m *Manager) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

// createTables creates the necessary database tables
func (m *Manager) createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			data TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS request_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL,
			method TEXT NOT NULL,
			headers TEXT,
			body TEXT,
			response TEXT,
			status INTEGER,
			duration_ms INTEGER,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS screenshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			url TEXT,
			file_path TEXT NOT NULL,
			size INTEGER,
			width INTEGER,
			height INTEGER,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_name ON sessions(name)`,
		`CREATE INDEX IF NOT EXISTS idx_request_logs_timestamp ON request_logs(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_screenshots_name ON screenshots(name)`,
	}
	
	for _, query := range queries {
		if _, err := m.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", query, err)
		}
	}
	
	return nil
}

// SaveSession saves a browser session
func (m *Manager) SaveSession(name string, sessionData *browser.SessionData) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"success": false,
	}
	
	dataJSON, err := json.Marshal(sessionData)
	if err != nil {
		result["message"] = fmt.Sprintf("Failed to serialize session data: %v", err)
		return result, nil
	}
	
	query := `INSERT OR REPLACE INTO sessions (name, data, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP)`
	_, err = m.db.Exec(query, name, string(dataJSON))
	if err != nil {
		result["message"] = fmt.Sprintf("Failed to save session: %v", err)
		return result, nil
	}
	
	result["success"] = true
	result["message"] = "Session saved successfully"
	result["name"] = name
	
	return result, nil
}

// LoadSession loads a browser session
func (m *Manager) LoadSession(name string) (*browser.SessionData, error) {
	query := `SELECT data FROM sessions WHERE name = ?`
	var dataJSON string
	
	err := m.db.QueryRow(query, name).Scan(&dataJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to load session: %w", err)
	}
	
	var sessionData browser.SessionData
	if err := json.Unmarshal([]byte(dataJSON), &sessionData); err != nil {
		return nil, fmt.Errorf("failed to deserialize session data: %w", err)
	}
	
	return &sessionData, nil
}

// ListSessions returns all saved sessions
func (m *Manager) ListSessions() ([]Session, error) {
	query := `SELECT id, name, data, created_at, updated_at FROM sessions ORDER BY updated_at DESC`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions: %w", err)
	}
	defer rows.Close()
	
	var sessions []Session
	for rows.Next() {
		var session Session
		var dataJSON string
		
		err := rows.Scan(&session.ID, &session.Name, &dataJSON, &session.CreatedAt, &session.UpdatedAt)
		if err != nil {
			continue
		}
		
		// Parse session data
		var sessionData browser.SessionData
		if json.Unmarshal([]byte(dataJSON), &sessionData) == nil {
			session.Data = &sessionData
		}
		
		sessions = append(sessions, session)
	}
	
	return sessions, nil
}

// DeleteSession deletes a saved session
func (m *Manager) DeleteSession(name string) error {
	query := `DELETE FROM sessions WHERE name = ?`
	result, err := m.db.Exec(query, name)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("session '%s' not found", name)
	}
	
	return nil
}

// LogRequest logs an HTTP request and response
func (m *Manager) LogRequest(url, method, headers, body, response string, status int, duration time.Duration) error {
	query := `INSERT INTO request_logs (url, method, headers, body, response, status, duration_ms) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	_, err := m.db.Exec(query, url, method, headers, body, response, status, duration.Milliseconds())
	if err != nil {
		return fmt.Errorf("failed to log request: %w", err)
	}
	
	return nil
}

// GetRequestLogs returns recent request logs
func (m *Manager) GetRequestLogs(limit int) ([]RequestLog, error) {
	query := `SELECT id, url, method, headers, body, response, status, duration_ms, timestamp 
			  FROM request_logs ORDER BY timestamp DESC LIMIT ?`
	
	rows, err := m.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query request logs: %w", err)
	}
	defer rows.Close()
	
	var logs []RequestLog
	for rows.Next() {
		var log RequestLog
		err := rows.Scan(&log.ID, &log.URL, &log.Method, &log.Headers, &log.Body, 
						&log.Response, &log.Status, &log.Duration, &log.Timestamp)
		if err != nil {
			continue
		}
		logs = append(logs, log)
	}
	
	return logs, nil
}

// SaveScreenshot saves screenshot metadata
func (m *Manager) SaveScreenshot(name, url, filePath, description string, size int64, width, height int) error {
	query := `INSERT INTO screenshots (name, url, file_path, size, width, height, description) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	_, err := m.db.Exec(query, name, url, filePath, size, width, height, description)
	if err != nil {
		return fmt.Errorf("failed to save screenshot metadata: %w", err)
	}
	
	return nil
}

// GetScreenshots returns screenshot metadata
func (m *Manager) GetScreenshots(limit int) ([]Screenshot, error) {
	query := `SELECT id, name, url, file_path, size, width, height, description, created_at 
			  FROM screenshots ORDER BY created_at DESC LIMIT ?`
	
	rows, err := m.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query screenshots: %w", err)
	}
	defer rows.Close()
	
	var screenshots []Screenshot
	for rows.Next() {
		var screenshot Screenshot
		err := rows.Scan(&screenshot.ID, &screenshot.Name, &screenshot.URL, &screenshot.FilePath,
						&screenshot.Size, &screenshot.Width, &screenshot.Height, 
						&screenshot.Description, &screenshot.CreatedAt)
		if err != nil {
			continue
		}
		screenshots = append(screenshots, screenshot)
	}
	
	return screenshots, nil
}

// CleanupOldData removes old data based on retention policies
func (m *Manager) CleanupOldData(maxAge time.Duration) error {
	cutoff := time.Now().Add(-maxAge)
	
	// Clean up old request logs
	query := `DELETE FROM request_logs WHERE timestamp < ?`
	if _, err := m.db.Exec(query, cutoff); err != nil {
		return fmt.Errorf("failed to cleanup request logs: %w", err)
	}
	
	// Clean up screenshot files older than maxAge
	screenshotQuery := `SELECT file_path FROM screenshots WHERE created_at < ?`
	rows, err := m.db.Query(screenshotQuery, cutoff)
	if err != nil {
		return fmt.Errorf("failed to query old screenshots: %w", err)
	}
	defer rows.Close()
	
	var filesToDelete []string
	for rows.Next() {
		var filePath string
		if rows.Scan(&filePath) == nil {
			filesToDelete = append(filesToDelete, filePath)
		}
	}
	
	// Delete screenshot files
	for _, filePath := range filesToDelete {
		os.Remove(filePath) // Ignore errors for cleanup
	}
	
	// Delete screenshot records
	deleteQuery := `DELETE FROM screenshots WHERE created_at < ?`
	if _, err := m.db.Exec(deleteQuery, cutoff); err != nil {
		return fmt.Errorf("failed to cleanup screenshot records: %w", err)
	}
	
	return nil
}

// GetStats returns database statistics
func (m *Manager) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Count sessions
	var sessionCount int
	m.db.QueryRow("SELECT COUNT(*) FROM sessions").Scan(&sessionCount)
	stats["sessions"] = sessionCount
	
	// Count request logs
	var requestCount int
	m.db.QueryRow("SELECT COUNT(*) FROM request_logs").Scan(&requestCount)
	stats["requests"] = requestCount
	
	// Count screenshots
	var screenshotCount int
	m.db.QueryRow("SELECT COUNT(*) FROM screenshots").Scan(&screenshotCount)
	stats["screenshots"] = screenshotCount
	
	// Database size
	fileInfo, err := os.Stat(m.path)
	if err == nil {
		stats["database_size"] = fileInfo.Size()
	}
	
	return stats, nil
}