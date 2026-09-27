package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/relentlessworks/unitkit/internal/model"
)

// MCPRequest represents a JSON-RPC 2.0 request.
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// MCPResponse represents a JSON-RPC 2.0 response.
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents a JSON-RPC 2.0 error.
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCPTool represents a tool definition.
type MCPTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema MCPSchema   `json:"inputSchema"`
}

// MCPSchema represents a JSON schema for tool input.
type MCPSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required,omitempty"`
}

func (h *Handler) mcp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed | hint: use POST /mcp with JSON-RPC 2.0 body")
		return
	}

	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusOK, MCPResponse{
			JSONRPC: "2.0",
			Error:   &MCPError{Code: -32700, Message: "Parse error: " + err.Error()},
		})
		return
	}

	resp := h.handleMCP(req)
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleMCP(req MCPRequest) MCPResponse {
	switch req.Method {
	case "initialize":
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]interface{}{
					"name":    "unitkit",
					"version": "0.1.0",
				},
			},
		}

	case "tools/list":
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": h.mcpTools(),
			},
		}

	case "tools/call":
		return h.handleToolCall(req)

	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &MCPError{Code: -32601, Message: "Method not found: " + req.Method},
		}
	}
}

func (h *Handler) mcpTools() []MCPTool {
	return []MCPTool{
		{
			Name:        "convert",
			Description: "Convert a value from one unit to another. Both units must be in the same category.",
			InputSchema: MCPSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"from": map[string]interface{}{
						"type":        "string",
						"description": "Source unit symbol (e.g. m, kg, C, L, Hz)",
					},
					"to": map[string]interface{}{
						"type":        "string",
						"description": "Target unit symbol (e.g. ft, lb, F, gal, kHz)",
					},
					"value": map[string]interface{}{
						"type":        "number",
						"description": "Value to convert",
					},
				},
				Required: []string{"from", "to", "value"},
			},
		},
		{
			Name:        "list_categories",
			Description: "List all unit categories (length, mass, temperature, volume, area, speed, time, data, pressure, energy, power, angle, frequency).",
			InputSchema: MCPSchema{
				Type:       "object",
				Properties: map[string]interface{}{},
			},
		},
		{
			Name:        "list_units",
			Description: "List all units in a specific category.",
			InputSchema: MCPSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"category": map[string]interface{}{
						"type":        "string",
						"description": "Category name (e.g. length, mass, temperature)",
					},
				},
				Required: []string{"category"},
			},
		},
	}
}

func (h *Handler) handleToolCall(req MCPRequest) MCPResponse {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &MCPError{Code: -32602, Message: "Invalid params: " + err.Error()},
		}
	}

	switch params.Name {
	case "convert":
		var args struct {
			From  string  `json:"from"`
			To    string  `json:"to"`
			Value float64 `json:"value"`
		}
		if err := json.Unmarshal(params.Arguments, &args); err != nil {
			return MCPResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &MCPError{Code: -32602, Message: "Invalid arguments: " + err.Error()},
			}
		}
		result, err := model.Convert(args.From, args.To, args.Value)
		if err != nil {
			return MCPResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &MCPError{Code: -32000, Message: err.Error()},
			}
		}
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": fmt.Sprintf("from=%s to=%s value=%s result=%s", result.From, result.To, model.FormatValue(result.Value), model.FormatValue(result.Result)),
					},
				},
			},
		}

	case "list_categories":
		cats := model.AllCategories()
		var sb strings.Builder
		for _, c := range cats {
			fmt.Fprintf(&sb, "name=%s description=%s base=%s units=%s\n", c.Name, c.Description, c.BaseUnit, strings.Join(c.Units, ","))
		}
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": sb.String(),
					},
				},
			},
		}

	case "list_units":
		var args struct {
			Category string `json:"category"`
		}
		if err := json.Unmarshal(params.Arguments, &args); err != nil {
			return MCPResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &MCPError{Code: -32602, Message: "Invalid arguments: " + err.Error()},
			}
		}
		units, ok := model.ListUnits(args.Category)
		if !ok {
			return MCPResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &MCPError{Code: -32000, Message: fmt.Sprintf("unknown category: %s", args.Category)},
			}
		}
		cat, _ := model.LookupCategory(args.Category)
		var sb strings.Builder
		fmt.Fprintf(&sb, "category=%s base=%s\n", cat.Name, cat.BaseUnit)
		for _, sym := range units {
			u, _ := model.LookupUnit(sym)
			fmt.Fprintf(&sb, "symbol=%s name=%s\n", u.Symbol, u.Name)
		}
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": sb.String(),
					},
				},
			},
		}

	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &MCPError{Code: -32601, Message: "Unknown tool: " + params.Name},
		}
	}
}
