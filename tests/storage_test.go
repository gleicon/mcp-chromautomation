package tests

import (
	"testing"
	"time"

	"github.com/gleicon/mcp-chromautomation/internal/browser"
	"github.com/gleicon/mcp-chromautomation/internal/storage"
)

func TestStorageManager(t *testing.T) {
	// Create storage manager 
	manager := storage.New()
	
	// Initialize
	if err := manager.Init(); err != nil {
		t.Fatalf("Failed to initialize storage manager: %v", err)
	}
	defer manager.Close()

	t.Run("SaveAndLoadSession", func(t *testing.T) {
		// Create test session data
		sessionData := &browser.SessionData{
			URL:       "https://example.com",
			UserAgent: "test-agent",
			Timestamp: time.Now(),
		}

		// Save session
		result, err := manager.SaveSession("test_session", sessionData)
		if err != nil {
			t.Fatalf("Failed to save session: %v", err)
		}

		if !result["success"].(bool) {
			t.Errorf("Expected success=true, got: %v", result)
		}

		// Load session
		loadedData, err := manager.LoadSession("test_session")
		if err != nil {
			t.Fatalf("Failed to load session: %v", err)
		}

		if loadedData.URL != sessionData.URL {
			t.Errorf("Expected URL %s, got %s", sessionData.URL, loadedData.URL)
		}

		if loadedData.UserAgent != sessionData.UserAgent {
			t.Errorf("Expected UserAgent %s, got %s", sessionData.UserAgent, loadedData.UserAgent)
		}
	})

	t.Run("ListSessions", func(t *testing.T) {
		sessions, err := manager.ListSessions()
		if err != nil {
			t.Fatalf("Failed to list sessions: %v", err)
		}

		// Should have at least the test session from previous test
		if len(sessions) == 0 {
			t.Error("Expected at least one session")
		}

		found := false
		for _, session := range sessions {
			if session.Name == "test_session" {
				found = true
				break
			}
		}

		if !found {
			t.Error("Expected to find test_session in session list")
		}
	})

	t.Run("LogRequest", func(t *testing.T) {
		err := manager.LogRequest(
			"https://example.com/api",
			"GET",
			`{"Authorization": "Bearer token"}`,
			"",
			`{"status": "ok"}`,
			200,
			150*time.Millisecond,
		)

		if err != nil {
			t.Fatalf("Failed to log request: %v", err)
		}

		// Get logs
		logs, err := manager.GetRequestLogs(10)
		if err != nil {
			t.Fatalf("Failed to get request logs: %v", err)
		}

		if len(logs) == 0 {
			t.Error("Expected at least one request log")
		}

		// Check the logged request
		found := false
		for _, log := range logs {
			if log.URL == "https://example.com/api" && log.Method == "GET" {
				found = true
				if log.Status != 200 {
					t.Errorf("Expected status 200, got %d", log.Status)
				}
				if log.Duration != 150 {
					t.Errorf("Expected duration 150ms, got %dms", log.Duration)
				}
				break
			}
		}

		if !found {
			t.Error("Expected to find logged request")
		}
	})

	t.Run("GetStats", func(t *testing.T) {
		stats, err := manager.GetStats()
		if err != nil {
			t.Fatalf("Failed to get stats: %v", err)
		}

		// Check that stats contain expected keys
		expectedKeys := []string{"sessions", "requests", "screenshots"}
		for _, key := range expectedKeys {
			if _, exists := stats[key]; !exists {
				t.Errorf("Expected stats to contain key: %s", key)
			}
		}

		// Sessions should be > 0 (from previous tests)
		if sessions, ok := stats["sessions"].(int); !ok || sessions == 0 {
			t.Errorf("Expected sessions > 0, got: %v", stats["sessions"])
		}

		// Requests should be > 0 (from previous tests)
		if requests, ok := stats["requests"].(int); !ok || requests == 0 {
			t.Errorf("Expected requests > 0, got: %v", stats["requests"])
		}
	})

	t.Run("DeleteSession", func(t *testing.T) {
		// Delete the test session
		err := manager.DeleteSession("test_session")
		if err != nil {
			t.Fatalf("Failed to delete session: %v", err)
		}

		// Try to load the deleted session (should fail)
		_, err = manager.LoadSession("test_session")
		if err == nil {
			t.Error("Expected error when loading deleted session")
		}

		// Try to delete non-existent session (should fail)
		err = manager.DeleteSession("non_existent_session")
		if err == nil {
			t.Error("Expected error when deleting non-existent session")
		}
	})
}