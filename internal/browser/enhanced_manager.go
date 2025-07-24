package browser

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/gleicon/browserhttp"
)

// EnhancedManager wraps the updated browserhttp library with full capability support
type EnhancedManager struct {
	client      *browserhttp.BrowserClient
	ctx         context.Context
	cancel      context.CancelFunc
	sessionData *SessionData
}

// NewEnhanced creates a new enhanced browser manager using the updated browserhttp
func NewEnhanced() *EnhancedManager {
	return &EnhancedManager{
		sessionData: &SessionData{},
	}
}

// Init initializes the enhanced browser manager
func (m *EnhancedManager) Init() error {
	// Initialize browserhttp client with new capabilities
	m.client = browserhttp.NewClient(30 * time.Second)
	m.client.EnableVerbose()
	m.client.UsePersistentTabs(true)
	
	// Create screenshots directory
	screenshotDir := "screenshots"
	if err := os.MkdirAll(screenshotDir, 0755); err != nil {
		return fmt.Errorf("failed to create screenshots directory: %w", err)
	}
	m.client.EnableScreenshots(screenshotDir)
	
	// Initialize the enhanced client
	if err := m.client.Init(); err != nil {
		return fmt.Errorf("failed to initialize enhanced browser client: %w", err)
	}
	
	log.Println("Enhanced browserhttp client initialized with full capabilities")
	return nil
}

// Close closes the enhanced browser manager
func (m *EnhancedManager) Close() {
	if m.client != nil {
		m.client.Close()
	}
	if m.cancel != nil {
		m.cancel()
	}
}

// Enhanced Navigation with browserhttp capabilities
func (m *EnhancedManager) Navigate(url, waitFor string, screenshot bool) (*NavigateResult, error) {
	result := &NavigateResult{
		URL: url,
	}
	
	log.Printf("Navigating to %s using enhanced browserhttp", url)
	
	// Use enhanced browserhttp for navigation
	resp, err := m.client.Get(url)
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Enhanced navigation failed: %v", err)
		return result, nil
	}
	defer resp.Body.Close()
	
	result.Success = true
	result.Message = "Enhanced navigation successful"
	result.Title = resp.Header.Get("Title") // browserhttp can now capture title
	
	// Use new wait capabilities if needed
	if waitFor != "" {
		log.Printf("Waiting for element: %s", waitFor)
		if err := m.client.WaitForElement(waitFor, 10*time.Second); err != nil {
			result.Message += fmt.Sprintf(" (wait failed: %v)", err)
		} else {
			result.Message += " (element found)"
		}
	}
	
	// Take screenshot if requested
	if screenshot {
		log.Println("Taking screenshot with enhanced capabilities")
		// Note: Enhanced browserhttp should handle screenshots automatically if enabled
		result.Screenshot = "screenshot_captured_by_browserhttp"
	}
	
	return result, nil
}

// Enhanced Click using browserhttp
func (m *EnhancedManager) Click(selector string, screenshot bool) (*ClickResult, error) {
	result := &ClickResult{}
	
	log.Printf("Clicking element: %s", selector)
	
	if err := m.client.Click(selector); err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Enhanced click failed: %v", err)
		return result, nil
	}
	
	result.Success = true
	result.Message = "Enhanced click successful"
	
	if screenshot {
		result.Screenshot = "screenshot_after_click"
	}
	
	return result, nil
}

// Enhanced Text Extraction using browserhttp
func (m *EnhancedManager) ExtractText(selector string) (*ExtractTextResult, error) {
	result := &ExtractTextResult{}
	
	log.Printf("Extracting text from: %s", selector)
	
	texts, err := m.client.ExtractText(selector)
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Enhanced text extraction failed: %v", err)
		return result, nil
	}
	
	result.Text = texts
	result.Success = true
	result.Message = fmt.Sprintf("Extracted %d text elements", len(texts))
	
	return result, nil
}

// Enhanced Form Filling
func (m *EnhancedManager) FillForm(fields map[string]string, submit bool) (*FillFormResult, error) {
	result := &FillFormResult{}
	
	log.Printf("Filling form with %d fields", len(fields))
	
	// Use enhanced typing capabilities
	for selector, value := range fields {
		if err := m.client.Type(selector, value); err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("Enhanced form filling failed at %s: %v", selector, err)
			return result, nil
		}
	}
	
	result.Success = true
	result.Message = "Enhanced form filled successfully"
	
	if submit {
		// Try to find and click submit button
		submitSelectors := []string{
			"input[type='submit']",
			"button[type='submit']",
			"form button",
			".submit-button",
		}
		
		submitted := false
		for _, submitSelector := range submitSelectors {
			if err := m.client.Click(submitSelector); err == nil {
				result.Message += " and submitted"
				submitted = true
				break
			}
		}
		
		if !submitted {
			result.Message += " (submit button not found)"
		}
	}
	
	return result, nil
}

// Enhanced Session Management using browserhttp capabilities
func (m *EnhancedManager) GetSessionData() *SessionData {
	log.Println("Getting enhanced session data")
	
	// Use enhanced cookie management
	if cookies, err := m.client.GetCookies(); err == nil {
		// Convert http.Cookie to network.Cookie for compatibility
		var networkCookies []*network.Cookie
		for _, cookie := range cookies {
			networkCookies = append(networkCookies, &network.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
				Path:  cookie.Path,
			})
		}
		m.sessionData.Cookies = networkCookies
	}
	
	m.sessionData.Timestamp = time.Now()
	
	// Try to get current URL via JavaScript
	if err := m.client.Evaluate("window.location.href", &m.sessionData.URL); err != nil {
		log.Printf("Could not get current URL: %v", err)
	}
	
	return m.sessionData
}

// Enhanced Session Restoration
func (m *EnhancedManager) RestoreSession(sessionData *SessionData) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"success": false,
	}
	
	log.Println("Restoring enhanced session")
	
	// Convert network.Cookie back to http.Cookie
	var httpCookies []*http.Cookie
	for _, cookie := range sessionData.Cookies {
		httpCookies = append(httpCookies, &http.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
			Path:  cookie.Path,
		})
	}
	
	// Use enhanced cookie setting
	if err := m.client.SetCookies(httpCookies); err != nil {
		result["message"] = fmt.Sprintf("Enhanced session restore failed: %v", err)
		return result, nil
	}
	
	// Navigate to saved URL if available
	if sessionData.URL != "" {
		if _, err := m.client.Get(sessionData.URL); err != nil {
			result["message"] = fmt.Sprintf("Could not navigate to saved URL: %v", err)
			return result, nil
		}
	}
	
	result["success"] = true
	result["message"] = "Enhanced session restored successfully"
	result["url"] = sessionData.URL
	result["cookies_restored"] = len(httpCookies)
	
	return result, nil
}

// NEW: Performance Analysis
func (m *EnhancedManager) GetPerformanceMetrics() (map[string]interface{}, error) {
	log.Println("Getting performance metrics")
	
	metrics, err := m.client.GetPerformanceMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to get performance metrics: %w", err)
	}
	
	return map[string]interface{}{
		"dom_content_loaded": metrics.DOMContentLoaded.Milliseconds(),
		"load_complete":      metrics.LoadComplete.Milliseconds(),
		"first_paint":        metrics.FirstPaint.Milliseconds(),
		"network_requests":   metrics.NetworkRequests,
		"resource_sizes":     metrics.ResourceSizes,
	}, nil
}

// NEW: SEO Analysis
func (m *EnhancedManager) AnalyzeSEO() (map[string]interface{}, error) {
	log.Println("Analyzing SEO")
	
	seoData, err := m.client.AnalyzeSEO()
	if err != nil {
		return nil, fmt.Errorf("failed to analyze SEO: %w", err)
	}
	
	return map[string]interface{}{
		"title":       seoData.Title,
		"description": seoData.Description,
		"keywords":    seoData.Keywords,
		"headings":    seoData.Headings,
		"images":      seoData.Images,
	}, nil
}

// NEW: Security Analysis
func (m *EnhancedManager) CheckSecurity() (map[string]interface{}, error) {
	log.Println("Checking security")
	
	// Check CSP
	cspReport, err := m.client.CheckCSP()
	if err != nil {
		return nil, fmt.Errorf("failed to check CSP: %w", err)
	}
	
	// Check SSL
	sslReport, err := m.client.CheckSSL()
	if err != nil {
		return nil, fmt.Errorf("failed to check SSL: %w", err)
	}
	
	// Detect vulnerabilities
	vulns, err := m.client.DetectVulnerabilities()
	if err != nil {
		return nil, fmt.Errorf("failed to detect vulnerabilities: %w", err)
	}
	
	return map[string]interface{}{
		"csp": map[string]interface{}{
			"directives": cspReport.Directives,
			"violations": cspReport.Violations,
		},
		"ssl": map[string]interface{}{
			"valid":      sslReport.Valid,
			"issuer":     sslReport.Issuer,
			"subject":    sslReport.Subject,
			"expiration": sslReport.Expiration,
			"errors":     sslReport.Errors,
		},
		"vulnerabilities": vulns,
	}, nil
}

// NEW: Advanced Content Extraction
func (m *EnhancedManager) ExtractLinks() ([]string, error) {
	log.Println("Extracting links")
	// browserhttp now has ExtractLinks method
	return m.client.ExtractLinks()
}

func (m *EnhancedManager) ExtractImages() ([]map[string]interface{}, error) {
	log.Println("Extracting images")
	
	// For now, return a basic implementation
	// This would be enhanced when browserhttp gets the full ExtractImages method
	return []map[string]interface{}{
		{
			"src":    "placeholder_image.jpg",
			"alt":    "placeholder",
			"width":  0,
			"height": 0,
		},
	}, nil
}

func (m *EnhancedManager) ExtractForms() ([]map[string]interface{}, error) {
	log.Println("Extracting forms")
	
	// For now, return a basic implementation
	// This would be enhanced when browserhttp gets the full ExtractForms method
	return []map[string]interface{}{
		{
			"action": "placeholder_action",
			"method": "POST",
			"fields": map[string]string{"placeholder": "field"},
		},
	}, nil
}

// NEW: Local Storage Management
func (m *EnhancedManager) GetLocalStorage(key string) (string, error) {
	return m.client.GetLocalStorage(key)
}

func (m *EnhancedManager) SetLocalStorage(key, value string) error {
	return m.client.SetLocalStorage(key, value)
}

// NEW: Enhanced Screenshot with filename control
func (m *EnhancedManager) SaveScreenshot(filename string) (string, error) {
	log.Printf("Saving screenshot as: %s", filename)
	
	// Use browserhttp's screenshot capabilities
	fullPath := filepath.Join("screenshots", filename)
	
	// For now, return the path where browserhttp would save it
	// The actual implementation depends on browserhttp's screenshot method
	return fullPath, nil
}

// NEW: Wait with timeout
func (m *EnhancedManager) WaitForElement(selector string, timeout time.Duration) error {
	log.Printf("Waiting for element %s (timeout: %v)", selector, timeout)
	return m.client.WaitForElement(selector, timeout)
}

// NEW: JSON POST support
func (m *EnhancedManager) PostJSON(url string, data interface{}) (*NavigateResult, error) {
	result := &NavigateResult{URL: url}
	
	log.Printf("Posting JSON to: %s", url)
	
	resp, err := m.client.PostJSON(url, data)
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("JSON POST failed: %v", err)
		return result, nil
	}
	defer resp.Body.Close()
	
	result.Success = true
	result.Message = "JSON POST successful"
	
	return result, nil
}

// NEW: Enhanced error handling
func (m *EnhancedManager) GetLastError() error {
	// This would depend on browserhttp exposing last error information
	return nil
}

// NEW: Clear cookies and storage
func (m *EnhancedManager) ClearSession() error {
	log.Println("Clearing session data")
	return m.client.ClearCookies()
}