package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestThinkToolServer(t *testing.T) {
	server := NewThinkToolServer("test-server")
	if server.serverName != "test-server" {
		t.Errorf("Expected server name to be 'test-server', got '%s'", server.serverName)
	}
}

func TestThinkTool(t *testing.T) {
	server := NewThinkToolServer("test-server")

	// Test think tool
	thought := "This is a test thought"
	ctx := context.Background()

	// Create tool request
	args := map[string]interface{}{
		"thought": thought,
	}

	request := mcp.CallToolRequest{}
	request.Params.Name = "think"
	request.Params.Arguments = args

	// Execute the think tool
	result, err := server.handleThink(ctx, request)
	if err != nil {
		t.Fatalf("Error calling think tool: %v", err)
	}

	// Get the text from the result
	if len(result.Content) == 0 {
		t.Fatalf("Expected content in result, got empty content")
	}

	textContent, ok := mcp.AsTextContent(result.Content[0])
	if !ok {
		t.Fatalf("Expected TextContent, got different type")
	}

	expectedText := "Thought recorded: This is a test thought"
	if textContent.Text != expectedText {
		t.Errorf("Expected text '%s', got '%s'", expectedText, textContent.Text)
	}

	// Check if thought was recorded
	if len(server.thoughts) != 1 {
		t.Fatalf("Expected 1 thought to be recorded, got %d", len(server.thoughts))
	}

	if server.thoughts[0].Thought != thought {
		t.Errorf("Expected thought '%s', got '%s'", thought, server.thoughts[0].Thought)
	}
}

func TestGetThoughtsAndStats(t *testing.T) {
	server := NewThinkToolServer("test-server")
	ctx := context.Background()

	// Add some thoughts
	thoughts := []string{
		"First thought",
		"Second thought that is longer",
		"Third thought",
	}

	for _, thought := range thoughts {
		args := map[string]interface{}{
			"thought": thought,
		}

		request := mcp.CallToolRequest{}
		request.Params.Name = "think"
		request.Params.Arguments = args

		_, err := server.handleThink(ctx, request)
		if err != nil {
			t.Fatalf("Error calling think tool: %v", err)
		}
	}

	// Test get_thoughts
	getRequest := mcp.CallToolRequest{}
	getRequest.Params.Name = "get_thoughts"
	getRequest.Params.Arguments = map[string]interface{}{}

	getResult, err := server.handleGetThoughts(ctx, getRequest)
	if err != nil {
		t.Fatalf("Error calling get_thoughts tool: %v", err)
	}

	// Check if the result contains all thoughts
	if len(getResult.Content) == 0 {
		t.Fatalf("Expected content in result, got empty content")
	}

	textContent, ok := mcp.AsTextContent(getResult.Content[0])
	if !ok {
		t.Fatalf("Expected TextContent, got different type")
	}

	for _, thought := range thoughts {
		if !strings.Contains(textContent.Text, thought) {
			t.Errorf("Expected get_thoughts result to contain '%s'", thought)
		}
	}

	// Test get_thought_stats
	statsRequest := mcp.CallToolRequest{}
	statsRequest.Params.Name = "get_thought_stats"
	statsRequest.Params.Arguments = map[string]interface{}{}

	statsResult, err := server.handleGetThoughtStats(ctx, statsRequest)
	if err != nil {
		t.Fatalf("Error calling get_thought_stats tool: %v", err)
	}

	if len(statsResult.Content) == 0 {
		t.Fatalf("Expected content in result, got empty content")
	}

	statsTextContent, ok := mcp.AsTextContent(statsResult.Content[0])
	if !ok {
		t.Fatalf("Expected TextContent, got different type")
	}

	// Parse the stats JSON
	var stats map[string]interface{}
	err = json.Unmarshal([]byte(statsTextContent.Text), &stats)
	if err != nil {
		t.Fatalf("Error parsing stats JSON: %v", err)
	}

	// Check stats
	if stats["total_thoughts"].(float64) != 3 {
		t.Errorf("Expected total_thoughts to be 3, got %v", stats["total_thoughts"])
	}

	if stats["longest_thought_index"].(float64) != 2 {
		t.Errorf("Expected longest_thought_index to be 2, got %v", stats["longest_thought_index"])
	}
}

func TestClearThoughts(t *testing.T) {
	server := NewThinkToolServer("test-server")
	ctx := context.Background()

	// Add a thought
	args := map[string]interface{}{
		"thought": "Thought to be cleared",
	}

	request := mcp.CallToolRequest{}
	request.Params.Name = "think"
	request.Params.Arguments = args

	_, err := server.handleThink(ctx, request)
	if err != nil {
		t.Fatalf("Error calling think tool: %v", err)
	}

	// Verify thought was added
	if len(server.thoughts) != 1 {
		t.Fatalf("Expected 1 thought to be recorded, got %d", len(server.thoughts))
	}

	// Clear thoughts
	clearRequest := mcp.CallToolRequest{}
	clearRequest.Params.Name = "clear_thoughts"
	clearRequest.Params.Arguments = map[string]interface{}{}

	clearResult, err := server.handleClearThoughts(ctx, clearRequest)
	if err != nil {
		t.Fatalf("Error calling clear_thoughts tool: %v", err)
	}

	if len(clearResult.Content) == 0 {
		t.Fatalf("Expected content in result, got empty content")
	}

	clearTextContent, ok := mcp.AsTextContent(clearResult.Content[0])
	if !ok {
		t.Fatalf("Expected TextContent, got different type")
	}

	expectedText := "Cleared 1 recorded thoughts."
	if clearTextContent.Text != expectedText {
		t.Errorf("Expected text '%s', got '%s'", expectedText, clearTextContent.Text)
	}

	// Verify thoughts were cleared
	if len(server.thoughts) != 0 {
		t.Errorf("Expected 0 thoughts after clearing, got %d", len(server.thoughts))
	}
}
