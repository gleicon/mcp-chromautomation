package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gleicon/mcp-chromautomation/internal/browser"
	"github.com/gleicon/mcp-chromautomation/internal/storage"
)

// Example demonstrating basic usage of the chrome automation service
func main() {
	// Initialize browser manager
	browserManager := browser.New()
	if err := browserManager.Init(); err != nil {
		log.Fatal("Failed to initialize browser:", err)
	}
	defer browserManager.Close()

	// Initialize storage manager
	storageManager := storage.New()
	if err := storageManager.Init(); err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}
	defer storageManager.Close()

	// Example 1: Basic navigation
	fmt.Println("🌐 Example 1: Basic Navigation")
	result, err := browserManager.Navigate("https://httpbin.org", "", true)
	if err != nil {
		log.Printf("Navigation failed: %v", err)
	} else {
		fmt.Printf("✅ Navigated to %s\n", result.URL)
		fmt.Printf("📄 Page title: %s\n", result.Title)
		if result.Screenshot != "" {
			fmt.Println("📸 Screenshot captured")
		}
	}

	// Example 2: Form interaction
	fmt.Println("\n📝 Example 2: Form Interaction")
	
	// Navigate to a form page
	formResult, err := browserManager.Navigate("https://httpbin.org/forms/post", "", false)
	if err != nil {
		log.Printf("Form navigation failed: %v", err)
	} else {
		fmt.Printf("✅ Navigated to form page: %s\n", formResult.URL)
		
		// Fill form fields
		fields := map[string]string{
			"input[name='custname']": "John Doe",
			"input[name='custtel']":  "+1234567890",
			"input[name='custemail']": "john@example.com",
			"select[name='size']":    "medium",
		}
		
		fillResult, err := browserManager.FillForm(fields, false)
		if err != nil {
			log.Printf("Form filling failed: %v", err)
		} else {
			fmt.Printf("✅ Form filled: %s\n", fillResult.Message)
		}
	}

	// Example 3: Text extraction
	fmt.Println("\n📄 Example 3: Text Extraction")
	
	// Navigate to a content page
	contentResult, err := browserManager.Navigate("https://httpbin.org/html", "", false)
	if err != nil {
		log.Printf("Content navigation failed: %v", err)
	} else {
		fmt.Printf("✅ Navigated to content page: %s\n", contentResult.URL)
		
		// Extract text from headings
		textResult, err := browserManager.ExtractText("h1")
		if err != nil {
			log.Printf("Text extraction failed: %v", err)
		} else {
			fmt.Printf("✅ Extracted %d headings:\n", len(textResult.Text))
			for i, text := range textResult.Text {
				fmt.Printf("  H1 #%d: %s\n", i+1, text)
			}
		}
	}

	// Example 4: Session management
	fmt.Println("\n💾 Example 4: Session Management")
	
	// Save current session
	sessionData := browserManager.GetSessionData()
	saveResult, err := storageManager.SaveSession("example_session", sessionData)
	if err != nil {
		log.Printf("Session save failed: %v", err)
	} else {
		fmt.Printf("✅ Session saved: %v\n", saveResult)
	}
	
	// List all sessions
	sessions, err := storageManager.ListSessions()
	if err != nil {
		log.Printf("Failed to list sessions: %v", err)
	} else {
		fmt.Printf("📋 Found %d saved sessions:\n", len(sessions))
		for _, session := range sessions {
			fmt.Printf("  - %s (URL: %s, Updated: %s)\n", 
				session.Name, 
				session.Data.URL, 
				session.UpdatedAt.Format("2006-01-02 15:04:05"))
		}
	}

	// Example 5: Request logging
	fmt.Println("\n📊 Example 5: Request Logging")
	
	// Log a sample request (normally done automatically by the browser manager)
	err = storageManager.LogRequest(
		"https://httpbin.org/get",
		"GET",
		`{"User-Agent": "mcp-chromautomation"}`,
		"",
		`{"args": {}, "headers": {...}}`,
		200,
		250*time.Millisecond,
	)
	if err != nil {
		log.Printf("Failed to log request: %v", err)
	} else {
		fmt.Println("✅ Request logged successfully")
	}
	
	// Get recent logs
	logs, err := storageManager.GetRequestLogs(5)
	if err != nil {
		log.Printf("Failed to get logs: %v", err)
	} else {
		fmt.Printf("📋 Recent %d requests:\n", len(logs))
		for _, reqLog := range logs {
			fmt.Printf("  %s %s - %d (%dms)\n", 
				reqLog.Method, 
				reqLog.URL, 
				reqLog.Status, 
				reqLog.Duration)
		}
	}

	// Example 6: Database statistics
	fmt.Println("\n📈 Example 6: Database Statistics")
	
	stats, err := storageManager.GetStats()
	if err != nil {
		log.Printf("Failed to get stats: %v", err)
	} else {
		fmt.Println("📊 Database Statistics:")
		for key, value := range stats {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	fmt.Println("\n🎉 All examples completed successfully!")
}