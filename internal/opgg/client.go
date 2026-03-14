package opgg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      int         `json:"id"`
}

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

type JSONRPCResponse struct {
	Result struct {
		Tools []Tool `json:"tools"`
	} `json:"result"`
}

func (c *Client) ListTools(ctx context.Context) ([]Tool, error) {
	payload := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/list",
		ID:      1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, err
	}

	return rpcResp.Result.Tools, nil
}

// CallToolRequest follows the MCP spec for invoking a tool
type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult represents the response from an MCP tool execution
type ToolResult struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	IsError bool `json:"isError,omitempty"`
}

// JSONRPCResponseCall matches the JSON-RPC envelope for a tool call result
type JSONRPCResponseCall struct {
	Result ToolResult `json:"result"`
}

// CallTool executes a specific tool on the OP.GG MCP server
func (c *Client) CallTool(ctx context.Context, name string, args map[string]any) (map[string]any, error) {
	// 1. Prepare the payload
	payload := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params: map[string]any{
			"name":      name,
			"arguments": args,
		},
		ID: 1,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// 2. Send the request
	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 3. Decode the raw JSON-RPC response first
	var rpcResp struct {
		Result map[string]any `json:"result"`
		Error  any            `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON-RPC: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("MCP server error: %v", rpcResp.Error)
	}

	// Return the 'Result' map directly to Gemini
	return rpcResp.Result, nil
}
