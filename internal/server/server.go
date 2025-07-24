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

// MCPChromeServer represents the enhanced MCP server using mark3labs/mcp-go
type MCPChromeServer struct {
	server  *server.MCPServer
	browser *browser.Manager
	storage *storage.Manager
}

// New creates a new MCP Chrome automation server
func New() *MCPChromeServer {
	// Create MCP server with proper configuration
	mcpServer := server.NewMCPServer(
		"mcp-chromautomation",
		"0.1.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
		server.WithLogging(),
	)

	// Initialize browser and storage managers
	browserManager := browser.New()
	storageManager := storage.New()

	chromeServer := &MCPChromeServer{
		server:  mcpServer,
		browser: browserManager,
		storage: storageManager,
	}

	// Register all tools
	chromeServer.registerTools()

	return chromeServer
}

// Start starts the MCP server
func (s *MCPChromeServer) Start(ctx context.Context) error {
	log.Println("Starting MCP Chrome Automation Server with mark3labs/mcp-go...")

	// Initialize components
	if err := s.storage.Init(); err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	if err := s.browser.Init(); err != nil {
		return fmt.Errorf("failed to initialize browser: %w", err)
	}

	defer func() {
		s.browser.Close()
		s.storage.Close()
	}()

	// Start the MCP server using stdio
	if err := server.ServeStdio(s.server); err != nil {
		return fmt.Errorf("failed to serve MCP server: %w", err)
	}

	return nil
}

// registerTools registers all Chrome automation tools with the MCP server
func (s *MCPChromeServer) registerTools() {
	// Chrome Navigate Tool
	navigateTool := mcp.NewTool("chrome_navigate",
		mcp.WithDescription("Navigate to a URL using Chrome browser with optional screenshot and element waiting"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The URL to navigate to")),
		mcp.WithString("wait_for",
			mcp.Description("CSS selector to wait for before completing navigation")),
		mcp.WithBoolean("screenshot",
			mcp.Description("Whether to take a screenshot after navigation"),
			mcp.DefaultBool(false)),
	)

	s.server.AddTool(navigateTool, s.handleChromeNavigate)

	// Chrome Click Tool
	clickTool := mcp.NewTool("chrome_click",
		mcp.WithDescription("Click an element on the current page"),
		mcp.WithString("selector",
			mcp.Required(),
			mcp.Description("CSS selector of the element to click")),
		mcp.WithBoolean("screenshot",
			mcp.Description("Whether to take a screenshot after clicking"),
			mcp.DefaultBool(false)),
	)

	s.server.AddTool(clickTool, s.handleChromeClick)

	// Chrome Extract Text Tool
	extractTextTool := mcp.NewTool("chrome_extract_text",
		mcp.WithDescription("Extract text content from elements on the current page"),
		mcp.WithString("selector",
			mcp.Required(),
			mcp.Description("CSS selector of elements to extract text from")),
	)

	s.server.AddTool(extractTextTool, s.handleChromeExtractText)

	// Chrome Fill Form Tool
	fillFormTool := mcp.NewTool("chrome_fill_form",
		mcp.WithDescription("Fill form fields on the current page"),
		mcp.WithObject("fields",
			mcp.Required(),
			mcp.Description("Map of CSS selectors to values for form fields")),
		mcp.WithBoolean("submit",
			mcp.Description("Whether to submit the form after filling"),
			mcp.DefaultBool(false)),
	)

	s.server.AddTool(fillFormTool, s.handleChromeFillForm)

	// Session Save Tool
	sessionSaveTool := mcp.NewTool("session_save",
		mcp.WithDescription("Save current browser session including cookies and URL for later restoration"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name for the saved session")),
	)

	s.server.AddTool(sessionSaveTool, s.handleSessionSave)

	// Session Load Tool
	sessionLoadTool := mcp.NewTool("session_load",
		mcp.WithDescription("Load a previously saved browser session"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the session to load")),
	)

	s.server.AddTool(sessionLoadTool, s.handleSessionLoad)

	// Enhanced Tools - Additional functionality
	screenshotTool := mcp.NewTool("chrome_screenshot",
		mcp.WithDescription("Take a screenshot of the current page"),
		mcp.WithString("filename",
			mcp.Description("Optional filename for the screenshot")),
	)

	s.server.AddTool(screenshotTool, s.handleChromeScreenshot)

	// Wait for Element Tool
	waitTool := mcp.NewTool("chrome_wait_for_element",
		mcp.WithDescription("Wait for an element to appear on the page"),
		mcp.WithString("selector",
			mcp.Required(),
			mcp.Description("CSS selector of the element to wait for")),
		mcp.WithNumber("timeout",
			mcp.Description("Timeout in seconds"),
			mcp.DefaultNumber(10)),
	)

	s.server.AddTool(waitTool, s.handleChromeWaitForElement)
}

// Tool handler implementations
func (s *MCPChromeServer) handleChromeNavigate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract arguments using the mcp helper functions
	url := mcp.ParseString(request, "url", "")
	if url == "" {
		return mcp.NewToolResultError("URL parameter is required"), nil
	}

	waitFor := mcp.ParseString(request, "wait_for", "")
	screenshot := mcp.ParseBoolean(request, "screenshot", false)

	// Execute navigation
	result, err := s.browser.Navigate(url, waitFor, screenshot)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Navigation failed", err), nil
	}

	// Convert result to JSON for MCP response
	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *MCPChromeServer) handleChromeClick(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector := mcp.ParseString(request, "selector", "")
	if selector == "" {
		return mcp.NewToolResultError("Selector parameter is required"), nil
	}

	screenshot := mcp.ParseBoolean(request, "screenshot", false)

	result, err := s.browser.Click(selector, screenshot)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Click failed", err), nil
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *MCPChromeServer) handleChromeExtractText(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector := mcp.ParseString(request, "selector", "")
	if selector == "" {
		return mcp.NewToolResultError("Selector parameter is required"), nil
	}

	result, err := s.browser.ExtractText(selector)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Text extraction failed", err), nil
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *MCPChromeServer) handleChromeFillForm(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	fieldsMap := mcp.ParseStringMap(request, "fields", map[string]any{})
	if len(fieldsMap) == 0 {
		return mcp.NewToolResultError("Fields parameter is required"), nil
	}

	// Convert map[string]any to map[string]string
	fields := make(map[string]string)
	for key, value := range fieldsMap {
		if strValue, ok := value.(string); ok {
			fields[key] = strValue
		}
	}

	submit := mcp.ParseBoolean(request, "submit", false)

	result, err := s.browser.FillForm(fields, submit)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Form filling failed", err), nil
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *MCPChromeServer) handleSessionSave(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := mcp.ParseString(request, "name", "")
	if name == "" {
		return mcp.NewToolResultError("Name parameter is required"), nil
	}

	sessionData := s.browser.GetSessionData()
	result, err := s.storage.SaveSession(name, sessionData)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Session save failed", err), nil
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *MCPChromeServer) handleSessionLoad(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := mcp.ParseString(request, "name", "")
	if name == "" {
		return mcp.NewToolResultError("Name parameter is required"), nil
	}

	sessionData, err := s.storage.LoadSession(name)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Session load failed", err), nil
	}

	result, err := s.browser.RestoreSession(sessionData)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Session restore failed", err), nil
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *MCPChromeServer) handleChromeScreenshot(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filename := mcp.ParseString(request, "filename", "")
	if filename == "" {
		filename = fmt.Sprintf("screenshot_%d.png", time.Now().Unix())
	}

	filePath, err := s.browser.SaveScreenshot(filename)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Screenshot failed", err), nil
	}

	result := map[string]interface{}{
		"success":   true,
		"message":   "Screenshot saved successfully",
		"file_path": filePath,
		"filename":  filename,
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func (s *MCPChromeServer) handleChromeWaitForElement(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	selector := mcp.ParseString(request, "selector", "")
	if selector == "" {
		return mcp.NewToolResultError("Selector parameter is required"), nil
	}

	timeout := mcp.ParseFloat64(request, "timeout", 10)
	if timeout <= 0 {
		timeout = 10 // Default 10 seconds
	}

	timeoutDuration := time.Duration(timeout * float64(time.Second))
	err := s.browser.WaitForElement(selector, timeoutDuration)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Wait failed", err), nil
	}

	result := map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Element '%s' found within %.0f seconds", selector, timeout),
	}

	resultJSON, _ := json.Marshal(result)
	return mcp.NewToolResultText(string(resultJSON)), nil
}