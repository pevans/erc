package a2

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/pevans/erc/mcp"
)

// StartMCPServer starts a new MCP service on localhost only, and listens on
// the configured port. It then starts a goroutine which forever-loops waiting
// for connections and runs a handler for them.
func (c *Computer) StartMCPServer() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", mcp.Port))
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go c.handleMCPConnection(conn)
		}
	}()

	return nil
}

// handleMCPConnection takes some connection and retrieves a request body from
// it, or waits until some request comes in. When a request does come in, if
// we can't understand it, we'll return an error to the requester. Otherwise,
// we'll run a request handler.
func (c *Computer) handleMCPConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	encoder := json.NewEncoder(conn)

	for scanner.Scan() {
		var req mcp.Request
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			encoder.Encode(mcp.Response{
				JSONRPC: "2.0",
				ID:      nil,
				Error:   &mcp.Error{Code: -32700, Message: "Parse error"},
			})
			continue
		}

		resp := c.handleMCPRequest(&req)
		encoder.Encode(resp)
	}
}

// handleMCPRequest returns a response for some given request. If we don't
// recognize the request, we'll return some error respoonse.
func (c *Computer) handleMCPRequest(req *mcp.Request) mcp.Response {
	switch req.Method {
	case "initialize":
		return mcp.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: mcp.InitializeResult{
				ProtocolVersion: "2024-11-05",
				ServerInfo: mcp.ServerInfo{
					Name:    "erc",
					Version: "0.1.0",
				},
				Capabilities: map[string]any{
					"tools": map[string]any{},
				},
			},
		}

	case "notifications/initialized":
		return mcp.Response{JSONRPC: "2.0", ID: req.ID}

	case "tools/list":
		return mcp.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: mcp.ToolsListResult{
				Tools: []mcp.Tool{
					{
						Name:        "status",
						Description: "Check if erc is running",
						InputSchema: mcp.InputSchema{Type: "object"},
					},
				},
			},
		}

	case "tools/call":
		var params mcp.ToolsCallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return mcp.Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &mcp.Error{Code: -32602, Message: "Invalid params"},
			}
		}

		if params.Name == "status" {
			return mcp.Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: mcp.ToolsCallResult{
					Content: []mcp.Content{
						{Type: "text", Text: "erc is running"},
					},
				},
			}
		}

		return mcp.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &mcp.Error{Code: -32601, Message: "Unknown tool"},
		}

	default:
		return mcp.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &mcp.Error{Code: -32601, Message: "Method not found"},
		}
	}
}
