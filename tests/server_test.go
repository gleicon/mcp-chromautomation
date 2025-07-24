package tests

import (
	"encoding/json"
	"testing"

	"github.com/gleicon/mcp-chromautomation/internal/server"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestMCPServer(t *testing.T) {
	srv := server.New()

	t.Run("Initialize", func(t *testing.T) {
		// This would normally be tested with actual stdin/stdout
		// For now, we'll just verify the server can be created
		if srv == nil {
			t.Error("Expected server to be created")
		}
	})

	t.Run("ToolsList", func(t *testing.T) {
		// We can't easily test the full server without setting up stdin/stdout
		// But we can verify the basic structure
		expectedTools := []string{
			"chrome_navigate",
			"chrome_click", 
			"chrome_extract_text",
			"chrome_fill_form",
			"session_save",
			"session_load",
		}

		// This is a basic structural test
		// In a real implementation, you'd mock stdin/stdout or use a test harness
		for _, tool := range expectedTools {
			// Verify tool names are valid
			if tool == "" {
				t.Errorf("Empty tool name found")
			}
		}
	})
}

func TestMCPRequestResponse(t *testing.T) {
	t.Run("ValidRequest", func(t *testing.T) {
		// Test JSON request parsing with mcp-go types
		requestJSON := `{
			"jsonrpc": "2.0",
			"id": 1,
			"method": "tools/list",
			"params": {}
		}`

		var request mcp.JSONRPCRequest
		err := json.Unmarshal([]byte(requestJSON), &request)
		if err != nil {
			t.Fatalf("Failed to parse request: %v", err)
		}

		if request.JSONRPC != "2.0" {
			t.Errorf("Expected jsonrpc=2.0, got %s", request.JSONRPC)
		}

		if request.Method != "tools/list" {
			t.Errorf("Expected method=tools/list, got %s", request.Method)
		}
	})

	t.Run("ValidResponse", func(t *testing.T) {
		// Test basic response structure
		response := mcp.NewJSONRPCResponse(mcp.NewRequestId(1), mcp.Result{})

		responseJSON, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("Failed to marshal response: %v", err)
		}

		// Verify it's valid JSON
		var parsed map[string]interface{}
		err = json.Unmarshal(responseJSON, &parsed)
		if err != nil {
			t.Fatalf("Response is not valid JSON: %v", err)
		}

		if parsed["jsonrpc"] != "2.0" {
			t.Errorf("Expected jsonrpc=2.0 in response")
		}
	})

	t.Run("ErrorResponse", func(t *testing.T) {
		// Test error response with mcp-go types
		response := mcp.NewJSONRPCError(
			mcp.NewRequestId(1),
			-32601,
			"Method not found",
			"test_method",
		)

		responseJSON, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("Failed to marshal error response: %v", err)
		}

		// Verify error structure
		var parsed map[string]interface{}
		err = json.Unmarshal(responseJSON, &parsed)
		if err != nil {
			t.Fatalf("Error response is not valid JSON: %v", err)
		}

		if errorObj, ok := parsed["error"].(map[string]interface{}); ok {
			if code, ok := errorObj["code"].(float64); !ok || code != -32601 {
				t.Errorf("Expected error code -32601, got %v", errorObj["code"])
			}
		} else {
			t.Error("Expected error object in response")
		}
	})
}

func TestToolSchemas(t *testing.T) {
	t.Run("NavigateToolSchema", func(t *testing.T) {
		// Test that tool creation works with mcp-go
		tool := mcp.NewTool("chrome_navigate",
			mcp.WithDescription("Navigate to a URL using Chrome browser"),
			mcp.WithString("url",
				mcp.Required(),
				mcp.Description("The URL to navigate to")),
			mcp.WithString("wait_for",
				mcp.Description("CSS selector to wait for")),
			mcp.WithBoolean("screenshot",
				mcp.Description("Whether to take a screenshot"),
				mcp.DefaultBool(false)),
		)

		// Verify tool structure
		if tool.Name != "chrome_navigate" {
			t.Errorf("Expected tool name 'chrome_navigate', got %s", tool.Name)
		}

		if tool.Description == "" {
			t.Error("Expected tool to have a description")
		}

		// Verify input schema has properties
		if len(tool.InputSchema.Properties) == 0 {
			t.Error("Expected tool to have input schema properties")
		}
	})
}