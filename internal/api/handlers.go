package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/relentlessworks/unitkit/internal/model"
)

// Handler holds the HTTP handler functions for the unitkit API.
type Handler struct {
	mux *http.ServeMux
}

// New creates a new API handler and registers all routes.
func New() *Handler {
	h := &Handler{mux: http.NewServeMux()}
	h.register()
	return h
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) register() {
	h.mux.HandleFunc("/help", h.help)
	h.mux.HandleFunc("/.well-known/agent.md", h.help)
	h.mux.HandleFunc("/categories", h.categories)
	h.mux.HandleFunc("/units/", h.units)
	h.mux.HandleFunc("/convert", h.convert)
	h.mux.HandleFunc("/convert/", h.convertPath)
	h.mux.HandleFunc("/mcp", h.mcp)
	h.mux.HandleFunc("/", h.root)
}

// wantsJSON checks if the client wants JSON responses.
func wantsJSON(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		return true
	}
	if r.URL.Query().Get("format") == "json" {
		return true
	}
	return false
}

// writeError writes an error response in plain text or JSON.
func writeError(w http.ResponseWriter, r *http.Request, status int, msg string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		// Try to split into error and hint
		parts := strings.SplitN(msg, " | hint: ", 2)
		errMsg := parts[0]
		hint := ""
		if len(parts) > 1 {
			hint = parts[1]
		}
		resp := map[string]string{"error": errMsg}
		if hint != "" {
			resp["hint"] = hint
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, "error: %s\n", msg)
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeText writes a plain text response.
func writeText(w http.ResponseWriter, status int, text string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprint(w, text)
}

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		writeError(w, r, http.StatusNotFound, fmt.Sprintf("unknown path: %s | hint: call GET /help for the operating manual", r.URL.Path))
		return
	}
	writeText(w, http.StatusOK, "unitkit — agentic-first unit conversion service\nCall GET /help for the operating manual\n")
}

func (h *Handler) help(w http.ResponseWriter, r *http.Request) {
	manual := `unitkit — Agentic-First Unit Conversion Service
==============================================

Convert between units of measurement across 13 categories: length, mass,
temperature, volume, area, speed, time, data, pressure, energy, power, angle,
and frequency. Pure stateless computation — no database, no auth required.

ENDPOINTS
---------

GET  /help                    This operating manual (also at /.well-known/agent.md)
GET  /categories              List all unit categories
GET  /units/{category}        List all units in a category
GET  /convert?from=X&to=Y&value=N   Convert N units of X to Y
POST /convert                 Convert via JSON body: {"from":"X","to":"Y","value":N}
GET  /convert/{from}/{to}/{value}  Convert via path params
POST /mcp                     MCP JSON-RPC 2.0 endpoint

RESPONSE FORMAT
---------------
Plain text by default (one labeled line per result).
JSON available via Accept: application/json header or ?format=json query param.

Plain text example:
  from=m to=km value=1000 result=1

JSON example:
  {"from":"m","to":"km","value":1000,"result":1,"category":"length"}

ERRORS
------
Errors include a hint for self-correction:
  error: unknown unit: foo | hint: call GET /categories to list all categories,
  then GET /units/{category} to see available units

CATEGORIES
----------
length, mass, temperature, volume, area, speed, time, data, pressure,
energy, power, angle, frequency

EXAMPLES
--------
curl "http://localhost:8080/convert?from=m&to=ft&value=1"
curl "http://localhost:8080/convert?from=C&to=F&value=25"
curl "http://localhost:8080/convert/m/ft/1"
curl -X POST http://localhost:8080/convert -d '{"from":"kg","to":"lb","value":1}'
curl "http://localhost:8080/units/length"
curl "http://localhost:8080/categories"

MCP
---
POST /mcp with JSON-RPC 2.0:
  {"jsonrpc":"2.0","method":"initialize","id":1}
  {"jsonrpc":"2.0","method":"tools/list","id":2}
  {"jsonrpc":"2.0","method":"tools/call","params":{"name":"convert","arguments":{"from":"m","to":"ft","value":1}},"id":3}
`
	writeText(w, http.StatusOK, manual)
}

func (h *Handler) categories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed | hint: use GET /categories")
		return
	}

	cats := model.AllCategories()

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, cats)
		return
	}

	var sb strings.Builder
	for _, c := range cats {
		fmt.Fprintf(&sb, "name=%s description=%s base=%s units=%s\n", c.Name, c.Description, c.BaseUnit, strings.Join(c.Units, ","))
	}
	writeText(w, http.StatusOK, sb.String())
}

func (h *Handler) units(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed | hint: use GET /units/{category}")
		return
	}

	// Extract category from path
	category := strings.TrimPrefix(r.URL.Path, "/units/")
	if category == "" {
		writeError(w, r, http.StatusBadRequest, "missing category | hint: call GET /categories to list all categories, then GET /units/{category}")
		return
	}

	units, ok := model.ListUnits(category)
	if !ok {
		writeError(w, r, http.StatusNotFound, fmt.Sprintf("unknown category: %s | hint: call GET /categories to list all categories", category))
		return
	}

	cat, _ := model.LookupCategory(category)

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, cat)
		return
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "category=%s base=%s\n", cat.Name, cat.BaseUnit)
	for _, sym := range units {
		u, _ := model.LookupUnit(sym)
		fmt.Fprintf(&sb, "symbol=%s name=%s\n", u.Symbol, u.Name)
	}
	writeText(w, http.StatusOK, sb.String())
}

func (h *Handler) convert(w http.ResponseWriter, r *http.Request) {
	var from, to string
	var value float64

	if r.Method == http.MethodGet {
		from = r.URL.Query().Get("from")
		to = r.URL.Query().Get("to")
		valStr := r.URL.Query().Get("value")
		if from == "" || to == "" || valStr == "" {
			writeError(w, r, http.StatusBadRequest, "missing required parameters: from, to, value | hint: use GET /convert?from=m&to=ft&value=1 or POST /convert with JSON body {\"from\":\"m\",\"to\":\"ft\",\"value\":1}")
			return
		}
		var err error
		value, err = parseFloat(valStr)
		if err != nil {
			writeError(w, r, http.StatusBadRequest, fmt.Sprintf("invalid value: %s | hint: value must be a number (e.g. 1, 3.14, -5)", valStr))
			return
		}
	} else if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, r, http.StatusBadRequest, "cannot read request body | hint: send a JSON body with from, to, and value fields")
			return
		}
		var req struct {
			From  string  `json:"from"`
			To    string  `json:"to"`
			Value float64 `json:"value"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			// Try form data
			from = r.FormValue("from")
			to = r.FormValue("to")
			valStr := r.FormValue("value")
			if from == "" || to == "" || valStr == "" {
				writeError(w, r, http.StatusBadRequest, "invalid request body | hint: send JSON {\"from\":\"m\",\"to\":\"ft\",\"value\":1} or form data with from, to, value")
				return
			}
			value, err = parseFloat(valStr)
			if err != nil {
				writeError(w, r, http.StatusBadRequest, fmt.Sprintf("invalid value: %s | hint: value must be a number", valStr))
				return
			}
		} else {
			from = req.From
			to = req.To
			value = req.Value
		}
		if from == "" || to == "" {
			writeError(w, r, http.StatusBadRequest, "missing from or to | hint: provide both 'from' and 'to' unit symbols")
			return
		}
	} else {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed | hint: use GET /convert?from=X&to=Y&value=N or POST /convert with JSON body")
		return
	}

	result, err := model.Convert(from, to, value)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, result)
		return
	}

	writeText(w, http.StatusOK, fmt.Sprintf("from=%s to=%s value=%s result=%s\n",
		result.From, result.To, model.FormatValue(result.Value), model.FormatValue(result.Result)))
}

func (h *Handler) convertPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, r, http.StatusMethodNotAllowed, "method not allowed | hint: use GET /convert/{from}/{to}/{value}")
		return
	}

	// Path: /convert/{from}/{to}/{value}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/convert/"), "/")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		writeError(w, r, http.StatusBadRequest, "invalid path format | hint: use GET /convert/{from}/{to}/{value} e.g. /convert/m/ft/1")
		return
	}

	from := parts[0]
	to := parts[1]
	valStr := parts[2]

	value, err := parseFloat(valStr)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, fmt.Sprintf("invalid value: %s | hint: value must be a number (e.g. 1, 3.14, -5)", valStr))
		return
	}

	result, err := model.Convert(from, to, value)
	if err != nil {
		writeError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, result)
		return
	}

	writeText(w, http.StatusOK, fmt.Sprintf("from=%s to=%s value=%s result=%s\n",
		result.From, result.To, model.FormatValue(result.Value), model.FormatValue(result.Result)))
}

// parseFloat parses a string to float64.
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%g", &f)
	return f, err
}
