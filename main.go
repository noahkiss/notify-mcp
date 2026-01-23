package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer("notify-mcp", "0.1.0",
		server.WithToolCapabilities(true),
	)

	tool := mcp.NewTool("send_notification",
		mcp.WithDescription("Send a push notification to the user's devices. Use when: long-running tasks complete, you need user input, or important errors occur."),
		mcp.WithString("title", mcp.Required(), mcp.Description("Short notification title")),
		mcp.WithString("message", mcp.Required(), mcp.Description("Notification body with details")),
	)

	s.AddTool(tool, handleNotification)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func handleNotification(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := req.RequireString("title")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	message, err := req.RequireString("message")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	appriseURL := os.Getenv("APPRISE_URL")
	if appriseURL == "" {
		return mcp.NewToolResultError("APPRISE_URL environment variable not set"), nil
	}

	tag := os.Getenv("APPRISE_TAG")
	if tag == "" {
		tag = "general"
	}

	payload := map[string]string{
		"tag":   tag,
		"title": title,
		"body":  message,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(appriseURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to send: %v", err)), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 || resp.StatusCode == 204 {
		return mcp.NewToolResultText("Notification sent successfully"), nil
	}
	return mcp.NewToolResultError(fmt.Sprintf("HTTP %d", resp.StatusCode)), nil
}
