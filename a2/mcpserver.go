package a2

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/pevans/erc/a2/a2state"
	"github.com/pevans/erc/elog"
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
	defer conn.Close() //nolint:errcheck

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var req mcp.Request
		if err := decoder.Decode(&req); err != nil {
			return
		}

		if resp := c.handleMCPRequest(&req); resp != nil {
			if err := encoder.Encode(resp); err != nil {
				return
			}
		}
	}
}

// handleMCPRequest returns a response for some given request. If we don't
// recognize the request, we'll return some error response. Returns nil for
// notifications (which don't expect responses).
func (c *Computer) handleMCPRequest(req *mcp.Request) *mcp.Response {
	switch req.Method {
	case "initialize":
		return &mcp.Response{
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
		return nil

	case "tools/list":
		return &mcp.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: mcp.ToolsListResult{
				Tools: []mcp.Tool{
					{
						Name:        "status",
						Description: "Check if erc is running",
						InputSchema: mcp.InputSchema{Type: "object"},
					},
					{
						Name:        "pause",
						Description: "Pause the emulator",
						InputSchema: mcp.InputSchema{Type: "object"},
					},
					{
						Name:        "resume",
						Description: "Resume the emulator",
						InputSchema: mcp.InputSchema{Type: "object"},
					},
					{
						Name:        "dbatch_start",
						Description: "Start a debug batch session",
						InputSchema: mcp.InputSchema{Type: "object"},
					},
					{
						Name:        "dbatch_stop",
						Description: "Stop a debug batch session",
						InputSchema: mcp.InputSchema{Type: "object"},
					},
				},
			},
		}

	case "tools/call":
		var params mcp.ToolsCallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return &mcp.Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &mcp.Error{Code: -32602, Message: "Invalid params"},
			}
		}

		return c.handleToolCall(req.ID, &params)

	default:
		return &mcp.Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &mcp.Error{Code: -32601, Message: "Method not found"},
		}
	}
}

func (c *Computer) handleToolCall(id any, params *mcp.ToolsCallParams) *mcp.Response {
	var text string

	switch params.Name {
	case "status":
		text = c.toolStatus()
	case "pause":
		text = c.toolPause()
	case "resume":
		text = c.toolResume()
	case "dbatch_start":
		text = c.toolDbatchStart()
	case "dbatch_stop":
		text = c.toolDbatchStop()
	default:
		return &mcp.Response{
			JSONRPC: "2.0",
			ID:      id,
			Error:   &mcp.Error{Code: -32601, Message: "Unknown tool"},
		}
	}

	return &mcp.Response{
		JSONRPC: "2.0",
		ID:      id,
		Result: mcp.ToolsCallResult{
			Content: []mcp.Content{
				{Type: "text", Text: text},
			},
		},
	}
}

func (c *Computer) toolStatus() string {
	return "erc is running"
}

func (c *Computer) toolPause() string {
	c.State.SetBool(a2state.Paused, true)
	return "emulator paused"
}

func (c *Computer) toolResume() string {
	c.State.SetBool(a2state.Paused, false)
	return "emulator resumed"
}

func (c *Computer) toolDbatchStart() string {
	c.dbatchMode = true
	c.dbatchTime = time.Now()
	c.instDiffMap = elog.NewInstructionMap()
	return "debug batch started"
}

func (c *Computer) toolDbatchStop() string {
	c.dbatchMode = false
	c.dbatchEnded = time.Now()

	if c.instDiffMap != nil && c.instDiffMapFileName != "" {
		_ = c.instDiffMap.WriteToFile(c.instDiffMapFileName)
		c.instDiffMap = nil
	}

	return "debug batch stopped"
}
