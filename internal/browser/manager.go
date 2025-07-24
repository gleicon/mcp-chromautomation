package browser

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/gleicon/browserhttp"
)

// Manager handles browser operations using the enhanced browserhttp library
type Manager struct {
	client      *browserhttp.BrowserClient
	ctx         context.Context
	cancel      context.CancelFunc
	sessionData *SessionData
}

// SessionData represents browser session information
type SessionData struct {
	Cookies   []*network.Cookie `json:"cookies"`
	URL       string            `json:"url"`
	UserAgent string            `json:"user_agent"`
	Timestamp time.Time         `json:"timestamp"`
}

// NavigateResult represents the result of a navigation operation
type NavigateResult struct {
	URL        string `json:"url"`
	Title      string `json:"title"`
	Screenshot string `json:"screenshot,omitempty"`
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
}

// ClickResult represents the result of a click operation
type ClickResult struct {
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
	Screenshot string `json:"screenshot,omitempty"`
}

// ExtractTextResult represents the result of text extraction
type ExtractTextResult struct {
	Text    []string `json:"text"`
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
}

// FillFormResult represents the result of form filling
type FillFormResult struct {
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
	Screenshot string `json:"screenshot,omitempty"`
}

// New creates a new browser manager
func New() *Manager {
	return &Manager{
		sessionData: &SessionData{},
	}
}

// Init initializes the browser manager to connect to existing Chrome
func (m *Manager) Init() error {
	// Try to connect to existing Chrome instance first
	ctx, cancel := m.connectToExistingChrome()
	if ctx != nil {
		m.ctx = ctx
		m.cancel = cancel
		log.Println("✅ Connected to existing Chrome instance")
	} else {
		// Fallback: create new context with remote debugging
		log.Println("⚠️  No existing Chrome found, will connect to Chrome with debugging enabled")
		log.Println("💡 Please start Chrome with: google-chrome --remote-debugging-port=9222")
		
		// Connect to Chrome on debugging port
		ctx, cancel := chromedp.NewRemoteAllocator(context.Background(), "http://localhost:9222")
		m.ctx, m.cancel = chromedp.NewContext(ctx)
		cancel() // We only need the allocator briefly
	}
	
	// Test the connection
	if err := m.testConnection(); err != nil {
		return fmt.Errorf("failed to connect to Chrome: %w\n\n💡 Make sure Chrome is running with debugging enabled:\n   google-chrome --remote-debugging-port=9222", err)
	}
	
	// Initialize browserhttp client with existing Chrome
	m.client = browserhttp.NewClient(30 * time.Second)
	m.client.EnableVerbose()
	m.client.UsePersistentTabs(true)
	
	// Create screenshots directory
	screenshotDir := "screenshots"
	if err := os.MkdirAll(screenshotDir, 0755); err != nil {
		return fmt.Errorf("failed to create screenshots directory: %w", err)
	}
	m.client.EnableScreenshots(screenshotDir)
	
	// Initialize browserhttp to use existing Chrome
	if err := m.initBrowserHTTPWithExistingChrome(); err != nil {
		log.Printf("⚠️  browserhttp fallback mode: %v", err)
	}
	
	return nil
}

// connectToExistingChrome attempts to connect to an existing Chrome instance
func (m *Manager) connectToExistingChrome() (context.Context, context.CancelFunc) {
	// Try common Chrome debugging ports
	ports := []string{"9222", "9223", "9224"}
	
	for _, port := range ports {
		debugURL := fmt.Sprintf("http://localhost:%s", port)
		
		// Test if Chrome debugging is available
		if m.isDebugPortAvailable(debugURL) {
			log.Printf("🔍 Found Chrome debugging on port %s", port)
			
			// Create remote allocator for existing Chrome
			allocCtx, allocCancel := chromedp.NewRemoteAllocator(context.Background(), debugURL)
			
			// Create new context connected to existing Chrome
			ctx, cancel := chromedp.NewContext(allocCtx)
			
			// Test the connection
			testCtx, testCancel := context.WithTimeout(ctx, 5*time.Second)
			defer testCancel()
			
			if err := chromedp.Run(testCtx, chromedp.Evaluate(`window.location.href`, nil)); err == nil {
				allocCancel() // Close allocator, keep the context
				return ctx, cancel
			}
			
			// Clean up failed attempt
			cancel()
			allocCancel()
		}
	}
	
	return nil, nil
}

// isDebugPortAvailable checks if Chrome debugging port is available
func (m *Manager) isDebugPortAvailable(debugURL string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(debugURL + "/json")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

// testConnection tests the Chrome connection
func (m *Manager) testConnection() error {
	testCtx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()
	
	// Try to get current URL to test connection
	var currentURL string
	return chromedp.Run(testCtx, chromedp.Evaluate(`window.location.href`, &currentURL))
}

// initBrowserHTTPWithExistingChrome configures browserhttp to use existing Chrome
func (m *Manager) initBrowserHTTPWithExistingChrome() error {
	// Note: browserhttp might need modifications to use existing Chrome
	// For now, we'll rely more on chromedp for operations
	log.Println("🔧 Configuring browserhttp to work with existing Chrome session")
	
	// This would require changes to your browserhttp library to accept
	// an existing Chrome debugging port instead of spawning new instance
	if err := m.client.Init(); err != nil {
		return fmt.Errorf("browserhttp init failed, using chromedp fallback: %w", err)
	}
	
	return nil
}

// Close closes the browser manager
func (m *Manager) Close() {
	if m.client != nil {
		m.client.Close()
	}
	if m.cancel != nil {
		m.cancel()
	}
}

// Navigate navigates to a URL using existing Chrome instance
func (m *Manager) Navigate(url, waitFor string, screenshot bool) (*NavigateResult, error) {
	result := &NavigateResult{
		URL: url,
	}
	
	log.Printf("🌐 Navigating to %s in existing Chrome instance", url)
	
	// Navigate using chromedp (which uses existing Chrome)
	var title string
	var screenshotData []byte
	
	actions := []chromedp.Action{
		chromedp.Navigate(url),
		chromedp.WaitReady("body"), // Wait for basic page load
		chromedp.Title(&title),
	}
	
	if waitFor != "" {
		log.Printf("⏳ Waiting for element: %s", waitFor)
		actions = append(actions, chromedp.WaitVisible(waitFor, chromedp.ByQuery))
	}
	
	if screenshot {
		log.Println("📸 Taking screenshot")
		actions = append(actions, chromedp.Screenshot("body", &screenshotData, chromedp.NodeVisible))
	}
	
	// Execute navigation with timeout
	navCtx, cancel := context.WithTimeout(m.ctx, 30*time.Second)
	defer cancel()
	
	if err := chromedp.Run(navCtx, actions...); err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Navigation failed: %v", err)
		return result, nil
	}
	
	result.Success = true
	result.Message = "Navigation successful using existing Chrome session"
	result.Title = title
	
	if screenshot && len(screenshotData) > 0 {
		result.Screenshot = base64.StdEncoding.EncodeToString(screenshotData)
		log.Printf("📸 Screenshot captured (%d bytes)", len(screenshotData))
	}
	
	// Also try browserhttp as a fallback/complement
	go func() {
		if resp, err := m.client.Get(url); err == nil {
			resp.Body.Close()
			log.Println("✅ browserhttp also confirmed navigation")
		}
	}()
	
	return result, nil
}

// Click clicks an element
func (m *Manager) Click(selector string, screenshot bool) (*ClickResult, error) {
	result := &ClickResult{}
	
	actions := []chromedp.Action{
		chromedp.Click(selector),
	}
	
	if screenshot {
		var screenshotData []byte
		actions = append(actions, chromedp.Screenshot("body", &screenshotData, chromedp.NodeVisible))
		
		if err := chromedp.Run(m.ctx, actions...); err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("Click failed: %v", err)
			return result, nil
		}
		
		if len(screenshotData) > 0 {
			result.Screenshot = base64.StdEncoding.EncodeToString(screenshotData)
		}
	} else {
		if err := chromedp.Run(m.ctx, actions...); err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("Click failed: %v", err)
			return result, nil
		}
	}
	
	result.Success = true
	result.Message = "Click successful"
	return result, nil
}

// ExtractText extracts text from elements
func (m *Manager) ExtractText(selector string) (*ExtractTextResult, error) {
	result := &ExtractTextResult{}
	
	var texts []string
	if err := chromedp.Run(m.ctx, chromedp.Evaluate(
		fmt.Sprintf(`Array.from(document.querySelectorAll('%s')).map(el => el.textContent.trim())`, selector),
		&texts,
	)); err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Text extraction failed: %v", err)
		return result, nil
	}
	
	result.Text = texts
	result.Success = true
	return result, nil
}

// FillForm fills form fields
func (m *Manager) FillForm(fields map[string]string, submit bool) (*FillFormResult, error) {
	result := &FillFormResult{}
	
	var actions []chromedp.Action
	
	// Fill each field
	for selector, value := range fields {
		actions = append(actions, chromedp.SendKeys(selector, value))
	}
	
	// Submit if requested
	if submit {
		// Look for common submit button selectors
		submitSelectors := []string{
			"input[type='submit']",
			"button[type='submit']",
			"form button",
			".submit-button",
			"[role='button'][type='submit']",
		}
		
		for _, submitSelector := range submitSelectors {
			actions = append(actions, chromedp.Click(submitSelector, chromedp.ByQuery))
			break // Use the first found submit button
		}
	}
	
	if err := chromedp.Run(m.ctx, actions...); err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Form filling failed: %v", err)
		return result, nil
	}
	
	result.Success = true
	result.Message = "Form filled successfully"
	if submit {
		result.Message += " and submitted"
	}
	
	return result, nil
}

// GetSessionData returns current session data
func (m *Manager) GetSessionData() *SessionData {
	// Get current cookies
	var cookies []*network.Cookie
	if err := chromedp.Run(m.ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		cookies, _ = network.GetCookies().Do(ctx)
		return nil
	})); err == nil {
		m.sessionData.Cookies = cookies
	}
	
	// Get current URL
	var currentURL string
	if err := chromedp.Run(m.ctx, chromedp.Location(&currentURL)); err == nil {
		m.sessionData.URL = currentURL
	}
	
	m.sessionData.Timestamp = time.Now()
	return m.sessionData
}

// RestoreSession restores session data
func (m *Manager) RestoreSession(sessionData *SessionData) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"success": false,
	}
	
	var actions []chromedp.Action
	
	// Navigate to the saved URL if available
	if sessionData.URL != "" {
		actions = append(actions, chromedp.Navigate(sessionData.URL))
	}
	
	// Restore cookies
	for _, cookie := range sessionData.Cookies {
		actions = append(actions, network.SetCookie(cookie.Name, cookie.Value).
			WithDomain(cookie.Domain).
			WithPath(cookie.Path).
			WithSecure(cookie.Secure).
			WithHTTPOnly(cookie.HTTPOnly))
	}
	
	if err := chromedp.Run(m.ctx, actions...); err != nil {
		result["message"] = fmt.Sprintf("Session restore failed: %v", err)
		return result, nil
	}
	
	result["success"] = true
	result["message"] = "Session restored successfully"
	result["url"] = sessionData.URL
	result["cookies_restored"] = len(sessionData.Cookies)
	
	return result, nil
}

// Enhanced methods using browserhttp improvements
func (m *Manager) PostWithJSON(url string, jsonData interface{}) (*NavigateResult, error) {
	// This is an enhancement we can add to browserhttp
	// For now, we'll use a placeholder implementation
	return m.Navigate(url, "", false)
}

func (m *Manager) GetWithCustomHeaders(url string, headers map[string]string) (*NavigateResult, error) {
	// Another enhancement for browserhttp
	return m.Navigate(url, "", false)
}

func (m *Manager) WaitForElement(selector string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()
	
	return chromedp.Run(ctx, chromedp.WaitVisible(selector))
}

func (m *Manager) SaveScreenshot(filename string) (string, error) {
	var screenshotData []byte
	
	if err := chromedp.Run(m.ctx, chromedp.Screenshot("body", &screenshotData, chromedp.NodeVisible)); err != nil {
		return "", fmt.Errorf("failed to take screenshot: %w", err)
	}
	
	fullPath := filepath.Join("screenshots", filename)
	if err := os.WriteFile(fullPath, screenshotData, 0644); err != nil {
		return "", fmt.Errorf("failed to save screenshot: %w", err)
	}
	
	return fullPath, nil
}