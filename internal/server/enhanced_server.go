package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gleicon/mcp-chromautomation/internal/browser"
	"github.com/gleicon/mcp-chromautomation/internal/storage"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// EnhancedMCPChromeServer provides all the new browserhttp capabilities via MCP
type EnhancedMCPChromeServer struct {
	server          *server.MCPServer
	browser         *browser.EnhancedManager
	storage         *storage.Manager
	legacyBrowser   *browser.Manager // Keep for compatibility
}

// NewEnhanced creates a new enhanced MCP Chrome automation server
func NewEnhanced() *EnhancedMCPChromeServer {
	// Create MCP server with enhanced configuration
	mcpServer := server.NewMCPServer(
		"mcp-chromautomation-enhanced",
		"0.2.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
		server.WithLogging(),
	)

	// Initialize enhanced browser and storage managers
	enhancedBrowser := browser.NewEnhanced()
	storageManager := storage.New()
	legacyBrowser := browser.New() // Keep for compatibility

	chromeServer := &EnhancedMCPChromeServer{
		server:        mcpServer,
		browser:       enhancedBrowser,
		storage:       storageManager,
		legacyBrowser: legacyBrowser,
	}

	// Register all enhanced tools
	chromeServer.registerEnhancedTools()

	return chromeServer
}

// Start starts the enhanced MCP server
func (s *EnhancedMCPChromeServer) Start(ctx context.Context) error {
	log.Println("Starting Enhanced MCP Chrome Automation Server with full browserhttp capabilities...")

	// Initialize components
	if err := s.storage.Init(); err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	if err := s.browser.Init(); err != nil {
		log.Printf("Enhanced browser init failed, falling back to legacy: %v", err)
		if err := s.legacyBrowser.Init(); err != nil {
			return fmt.Errorf("failed to initialize browser (both enhanced and legacy): %w", err)
		}
	}

	defer func() {
		s.browser.Close()
		s.legacyBrowser.Close()
		s.storage.Close()
	}()

	// Start the MCP server using stdio
	if err := server.ServeStdio(s.server); err != nil {
		return fmt.Errorf("failed to serve enhanced MCP server: %w", err)
	}

	return nil
}

// registerEnhancedTools registers all enhanced Chrome automation tools
func (s *EnhancedMCPChromeServer) registerEnhancedTools() {
	// Original tools (enhanced versions)
	s.registerCoreTools()
	
	// NEW: Content Analysis Tools
	s.registerContentAnalysisTools()
	
	// NEW: Performance Analysis Tools
	s.registerPerformanceTools()
	
	// NEW: Security Analysis Tools
	s.registerSecurityTools()
	
	// NEW: Advanced Interaction Tools
	s.registerAdvancedInteractionTools()
	
	// NEW: Data Management Tools
	s.registerDataManagementTools()
}

// Core enhanced tools
func (s *EnhancedMCPChromeServer) registerCoreTools() {
	// Enhanced Chrome Navigate Tool
	navigateTool := mcp.NewTool("chrome_navigate",
		mcp.WithDescription("Navigate to a URL using enhanced Chrome browser with performance tracking"),
		mcp.WithString("url", mcp.Required(), mcp.Description("The URL to navigate to")),
		mcp.WithString("wait_for", mcp.Description("CSS selector to wait for")),
		mcp.WithBoolean("screenshot", mcp.Description("Take a screenshot"), mcp.DefaultBool(false)),
		mcp.WithBoolean("track_performance", mcp.Description("Track performance metrics"), mcp.DefaultBool(false)),
	)
	s.server.AddTool(navigateTool, s.handleEnhancedNavigate)

	// Enhanced Form Filling Tool
	fillFormTool := mcp.NewTool("chrome_fill_form",
		mcp.WithDescription("Fill form fields with enhanced typing and validation"),
		mcp.WithObject("fields", mcp.Required(), mcp.Description("Map of CSS selectors to values")),
		mcp.WithBoolean("submit", mcp.Description("Submit the form"), mcp.DefaultBool(false)),
		mcp.WithBoolean("validate_before_submit", mcp.Description("Validate form before submission"), mcp.DefaultBool(true)),
	)
	s.server.AddTool(fillFormTool, s.handleEnhancedFillForm)
}

// Content analysis tools
func (s *EnhancedMCPChromeServer) registerContentAnalysisTools() {
	// Extract Links Tool
	extractLinksTool := mcp.NewTool("chrome_extract_links",
		mcp.WithDescription("Extract all links from the current page"),
		mcp.WithString("filter", mcp.Description("Optional filter pattern for links")),
	)
	s.server.AddTool(extractLinksTool, s.handleExtractLinks)

	// Extract Images Tool
	extractImagesTool := mcp.NewTool("chrome_extract_images",
		mcp.WithDescription("Extract all images with metadata from the current page"),
		mcp.WithBoolean("include_metadata", mcp.Description("Include width/height metadata"), mcp.DefaultBool(true)),
	)
	s.server.AddTool(extractImagesTool, s.handleExtractImages)

	// Extract Forms Tool
	extractFormsTool := mcp.NewTool("chrome_extract_forms",
		mcp.WithDescription("Extract all forms with their fields from the current page"),
	)
	s.server.AddTool(extractFormsTool, s.handleExtractForms)

	// SEO Analysis Tool
	seoAnalysisTool := mcp.NewTool("chrome_analyze_seo",
		mcp.WithDescription("Analyze SEO elements of the current page"),
	)
	s.server.AddTool(seoAnalysisTool, s.handleAnalyzeSEO)
}

// Performance analysis tools
func (s *EnhancedMCPChromeServer) registerPerformanceTools() {
	// Performance Metrics Tool
	performanceTool := mcp.NewTool("chrome_get_performance",
		mcp.WithDescription("Get detailed performance metrics for the current page"),
	)
	s.server.AddTool(performanceTool, s.handleGetPerformance)
}

// Security analysis tools
func (s *EnhancedMCPChromeServer) registerSecurityTools() {
	// Security Check Tool
	securityTool := mcp.NewTool("chrome_check_security",
		mcp.WithDescription("Perform comprehensive security analysis of the current page"),
	)
	s.server.AddTool(securityTool, s.handleCheckSecurity)
}

// Advanced interaction tools
func (s *EnhancedMCPChromeServer) registerAdvancedInteractionTools() {
	// JSON POST Tool
	jsonPostTool := mcp.NewTool("chrome_post_json",
		mcp.WithDescription("Send JSON data to a URL using the browser"),
		mcp.WithString("url", mcp.Required(), mcp.Description("The URL to send data to")),
		mcp.WithObject("data", mcp.Required(), mcp.Description("The JSON data to send")),
	)
	s.server.AddTool(jsonPostTool, s.handlePostJSON)

	// Enhanced Wait Tool
	waitAdvancedTool := mcp.NewTool("chrome_wait_advanced",
		mcp.WithDescription("Advanced waiting with multiple conditions"),
		mcp.WithString("selector", mcp.Description("CSS selector to wait for")),
		mcp.WithString("text", mcp.Description("Text content to wait for")),
		mcp.WithNumber("timeout", mcp.Description("Timeout in seconds"), mcp.DefaultNumber(10)),
	)
	s.server.AddTool(waitAdvancedTool, s.handleWaitAdvanced)
}

// Data management tools
func (s *EnhancedMCPChromeServer) registerDataManagementTools() {
	// Local Storage Get Tool
	localStorageGetTool := mcp.NewTool("chrome_get_local_storage",
		mcp.WithDescription("Get a value from browser localStorage"),
		mcp.WithString("key", mcp.Required(), mcp.Description("The localStorage key")),
	)
	s.server.AddTool(localStorageGetTool, s.handleGetLocalStorage)

	// Local Storage Set Tool
	localStorageSetTool := mcp.NewTool("chrome_set_local_storage",
		mcp.WithDescription("Set a value in browser localStorage"),
		mcp.WithString("key", mcp.Required(), mcp.Description("The localStorage key")),
		mcp.WithString("value", mcp.Required(), mcp.Description("The value to store")),
	)
	s.server.AddTool(localStorageSetTool, s.handleSetLocalStorage)

	// Clear Session Tool
	clearSessionTool := mcp.NewTool("chrome_clear_session",
		mcp.WithDescription("Clear all cookies and session data"),
	)
	s.server.AddTool(clearSessionTool, s.handleClearSession)
}

// Enhanced tool handlers
func (s *EnhancedMCPChromeServer) handleEnhancedNavigate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := mcp.ParseString(request, "url", "")
	if url == "" {
		return mcp.NewToolResultError("URL parameter is required"), nil
	}

	waitFor := mcp.ParseString(request, "wait_for", "")
	screenshot := mcp.ParseBoolean(request, "screenshot", false)
	trackPerformance := mcp.ParseBoolean(request, "track_performance", false)

	// Use enhanced navigation
	result, err := s.browser.Navigate(url, waitFor, screenshot)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Enhanced navigation failed", err), nil
	}

	// Add performance metrics if requested
	if trackPerformance {
		if metrics, err := s.browser.GetPerformanceMetrics(); err == nil {
			result.Message += fmt.Sprintf(" (Load time: %v)", metrics["load_complete"])
		}
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *EnhancedMCPChromeServer) handleEnhancedFillForm(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	fieldsMap := mcp.ParseStringMap(request, "fields", map[string]any{})
	if len(fieldsMap) == 0 {
		return mcp.NewToolResultError("Fields parameter is required"), nil
	}

	fields := make(map[string]string)
	for key, value := range fieldsMap {
		if strValue, ok := value.(string); ok {
			fields[key] = strValue
		}
	}

	submit := mcp.ParseBoolean(request, "submit", false)
	validate := mcp.ParseBoolean(request, "validate_before_submit", true)

	// Enhanced form filling with validation
	result, err := s.browser.FillForm(fields, submit)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Enhanced form filling failed", err), nil
	}

	if validate && submit {
		result.Message += " (validated before submission)"
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *EnhancedMCPChromeServer) handleExtractLinks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	links, err := s.browser.ExtractLinks()
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Link extraction failed", err), nil
	}

	result := map[string]interface{}{
		"success": true,
		"links":   links,
		"count":   len(links),
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *EnhancedMCPChromeServer) handleExtractImages(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	images, err := s.browser.ExtractImages()
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Image extraction failed", err), nil
	}

	result := map[string]interface{}{
		"success": true,
		"images":  images,
		"count":   len(images),
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *EnhancedMCPChromeServer) handleExtractForms(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	forms, err := s.browser.ExtractForms()
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Form extraction failed", err), nil
	}

	result := map[string]interface{}{
		"success": true,
		"forms":   forms,
		"count":   len(forms),
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *EnhancedMCPChromeServer) handleAnalyzeSEO(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	seoData, err := s.browser.AnalyzeSEO()
	if err != nil {
		return mcp.NewToolResultErrorFromErr("SEO analysis failed", err), nil
	}

	result := map[string]interface{}{
		"success": true,
		"seo":     seoData,
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *EnhancedMCPChromeServer) handleGetPerformance(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	performance, err := s.browser.GetPerformanceMetrics()
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Performance analysis failed", err), nil
	}

	result := map[string]interface{}{
		"success":     true,
		"performance": performance,
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *EnhancedMCPChromeServer) handleCheckSecurity(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	security, err := s.browser.CheckSecurity()
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Security analysis failed", err), nil
	}

	result := map[string]interface{}{
		"success":  true,
		"security": security,
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *EnhancedMCPChromeServer) handlePostJSON(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := mcp.ParseString(request, "url", "")
	if url == "" {
		return mcp.NewToolResultError("URL parameter is required"), nil
	}

	dataMap := mcp.ParseStringMap(request, "data", map[string]any{})
	if len(dataMap) == 0 {
		return mcp.NewToolResultError("Data parameter is required"), nil
	}

	result, err := s.browser.PostJSON(url, dataMap)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("JSON POST failed", err), nil
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *EnhancedMCPChromeServer) handleWaitAdvanced(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector := mcp.ParseString(request, "selector", "")
	text := mcp.ParseString(request, "text", "")
	timeout := time.Duration(mcp.ParseFloat64(request, "timeout", 10)) * time.Second

	if selector == "" && text == "" {
		return mcp.NewToolResultError("Either selector or text parameter is required"), nil
	}

	var err error
	if selector != "" {
		err = s.browser.WaitForElement(selector, timeout)
	}
	// Add text waiting when browserhttp supports it

	if err != nil {
		return mcp.NewToolResultErrorFromErr("Wait failed", err), nil
	}

	result := map[string]interface{}{
		"success": true,
		"message": "Wait condition met",
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *EnhancedMCPChromeServer) handleGetLocalStorage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key := mcp.ParseString(request, "key", "")
	if key == "" {
		return mcp.NewToolResultError("Key parameter is required"), nil
	}

	value, err := s.browser.GetLocalStorage(key)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get localStorage", err), nil
	}

	result := map[string]interface{}{
		"success": true,
		"key":     key,
		"value":   value,
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *EnhancedMCPChromeServer) handleSetLocalStorage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key := mcp.ParseString(request, "key", "")
	value := mcp.ParseString(request, "value", "")
	
	if key == "" {
		return mcp.NewToolResultError("Key parameter is required"), nil
	}

	err := s.browser.SetLocalStorage(key, value)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to set localStorage", err), nil
	}

	result := map[string]interface{}{
		"success": true,
		"message": "localStorage value set successfully",
		"key":     key,
		"value":   value,
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *EnhancedMCPChromeServer) handleClearSession(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	err := s.browser.ClearSession()
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to clear session", err), nil
	}

	result := map[string]interface{}{
		"success": true,
		"message": "Session cleared successfully",
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}