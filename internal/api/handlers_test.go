package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoot(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("root status = %d, want %d", w.Code, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "unitkit") {
		t.Errorf("root body should contain 'unitkit', got: %s", w.Body.String())
	}
}

func TestNotFound(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("unknown path status = %d, want %d", w.Code, http.StatusNotFound)
	}
	if !strings.Contains(w.Body.String(), "hint:") {
		t.Errorf("not found response should contain hint, got: %s", w.Body.String())
	}
}

func TestHelp(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/help", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("help status = %d, want %d", w.Code, http.StatusOK)
	}
	body := w.Body.String()
	if !strings.Contains(body, "unitkit") {
		t.Error("help should mention unitkit")
	}
	if !strings.Contains(body, "ENDPOINTS") {
		t.Error("help should list endpoints")
	}
	if !strings.Contains(body, "/convert") {
		t.Error("help should mention /convert")
	}
	if !strings.Contains(body, "/mcp") {
		t.Error("help should mention /mcp")
	}
}

func TestHelpWellKnown(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/.well-known/agent.md", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("agent.md status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestCategories(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("categories status = %d, want %d", w.Code, http.StatusOK)
	}
	body := w.Body.String()
	for _, cat := range []string{"length", "mass", "temperature", "volume", "area", "speed", "time", "data", "pressure", "energy", "power", "angle", "frequency"} {
		if !strings.Contains(body, cat) {
			t.Errorf("categories should contain %q, got: %s", cat, body)
		}
	}
}

func TestCategoriesJSON(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/categories?format=json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("categories JSON status = %d, want %d", w.Code, http.StatusOK)
	}
	var cats []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &cats); err != nil {
		t.Errorf("categories JSON parse error: %v, body: %s", err, w.Body.String())
	}
	if len(cats) != 13 {
		t.Errorf("categories JSON returned %d categories, want 13", len(cats))
	}
}

func TestUnits(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/units/length", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("units status = %d, want %d", w.Code, http.StatusOK)
	}
	body := w.Body.String()
	if !strings.Contains(body, "category=length") {
		t.Errorf("units should show category, got: %s", body)
	}
	if !strings.Contains(body, "symbol=m") {
		t.Errorf("units should list meter, got: %s", body)
	}
}

func TestUnitsUnknownCategory(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/units/nonexistent", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("unknown category status = %d, want %d", w.Code, http.StatusNotFound)
	}
	if !strings.Contains(w.Body.String(), "hint:") {
		t.Errorf("error should contain hint, got: %s", w.Body.String())
	}
}

func TestConvertGET(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/convert?from=m&to=ft&value=1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("convert GET status = %d, want %d", w.Code, http.StatusOK)
	}
	body := w.Body.String()
	if !strings.Contains(body, "from=m") {
		t.Errorf("convert response should contain from=m, got: %s", body)
	}
	if !strings.Contains(body, "to=ft") {
		t.Errorf("convert response should contain to=ft, got: %s", body)
	}
	if !strings.Contains(body, "result=") {
		t.Errorf("convert response should contain result=, got: %s", body)
	}
}

func TestConvertGETJSON(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/convert?from=C&to=F&value=25&format=json", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("convert JSON status = %d, want %d", w.Code, http.StatusOK)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("convert JSON parse error: %v, body: %s", err, w.Body.String())
	}
	if result["from"] != "C" {
		t.Errorf("convert JSON from = %v, want C", result["from"])
	}
	if result["to"] != "F" {
		t.Errorf("convert JSON to = %v, want F", result["to"])
	}
	if result["result"].(float64) != 77 {
		t.Errorf("convert JSON result = %v, want 77", result["result"])
	}
}

func TestConvertPOST(t *testing.T) {
	h := New()
	body := strings.NewReader(`{"from":"kg","to":"lb","value":1}`)
	req := httptest.NewRequest("POST", "/convert", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("convert POST status = %d, want %d", w.Code, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "from=kg") {
		t.Errorf("convert POST should contain from=kg, got: %s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "to=lb") {
		t.Errorf("convert POST should contain to=lb, got: %s", w.Body.String())
	}
}

func TestConvertPath(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/convert/m/km/1000", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("convert path status = %d, want %d", w.Code, http.StatusOK)
	}
	if !strings.Contains(w.Body.String(), "result=1") {
		t.Errorf("convert path should show result=1, got: %s", w.Body.String())
	}
}

func TestConvertMissingParams(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/convert?from=m&to=ft", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("missing params status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	if !strings.Contains(w.Body.String(), "hint:") {
		t.Errorf("error should contain hint, got: %s", w.Body.String())
	}
}

func TestConvertUnknownUnit(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/convert?from=foo&to=bar&value=1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("unknown unit status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	if !strings.Contains(w.Body.String(), "hint:") {
		t.Errorf("error should contain hint, got: %s", w.Body.String())
	}
}

func TestConvertCrossCategory(t *testing.T) {
	h := New()
	req := httptest.NewRequest("GET", "/convert?from=m&to=kg&value=1", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("cross-category status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	if !strings.Contains(w.Body.String(), "different categories") {
		t.Errorf("error should mention different categories, got: %s", w.Body.String())
	}
}

func TestMCPInitialize(t *testing.T) {
	h := New()
	body := strings.NewReader(`{"jsonrpc":"2.0","method":"initialize","id":1}`)
	req := httptest.NewRequest("POST", "/mcp", body)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("MCP initialize status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("MCP initialize parse error: %v, body: %s", err, w.Body.String())
	}
	if resp["jsonrpc"] != "2.0" {
		t.Errorf("MCP jsonrpc = %v, want 2.0", resp["jsonrpc"])
	}
}

func TestMCPToolsList(t *testing.T) {
	h := New()
	body := strings.NewReader(`{"jsonrpc":"2.0","method":"tools/list","id":2}`)
	req := httptest.NewRequest("POST", "/mcp", body)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("MCP tools/list status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("MCP tools/list parse error: %v, body: %s", err, w.Body.String())
	}
	tools, ok := resp["result"].(map[string]interface{})["tools"].([]interface{})
	if !ok {
		t.Fatal("MCP tools/list should return tools array")
	}
	if len(tools) != 3 {
		t.Errorf("MCP tools/list returned %d tools, want 3", len(tools))
	}
}

func TestMCPToolsCallConvert(t *testing.T) {
	h := New()
	body := strings.NewReader(`{"jsonrpc":"2.0","method":"tools/call","params":{"name":"convert","arguments":{"from":"m","to":"ft","value":1}},"id":3}`)
	req := httptest.NewRequest("POST", "/mcp", body)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("MCP tools/call status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("MCP tools/call parse error: %v, body: %s", err, w.Body.String())
	}
	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		t.Fatal("MCP tools/call should return result")
	}
	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatal("MCP tools/call should return content array")
	}
	textItem, ok := content[0].(map[string]interface{})
	if !ok {
		t.Fatal("MCP tools/call content[0] should be a map")
	}
	if textItem["type"] != "text" {
		t.Errorf("MCP content type = %v, want text", textItem["type"])
	}
	text, ok := textItem["text"].(string)
	if !ok || !strings.Contains(text, "from=m") {
		t.Errorf("MCP text should contain from=m, got: %v", textItem["text"])
	}
}

func TestMCPToolsCallListCategories(t *testing.T) {
	h := New()
	body := strings.NewReader(`{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_categories","arguments":{}},"id":4}`)
	req := httptest.NewRequest("POST", "/mcp", body)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("MCP list_categories status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("MCP list_categories parse error: %v, body: %s", err, w.Body.String())
	}
	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		t.Fatal("MCP list_categories should return result")
	}
	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatal("MCP list_categories should return content array")
	}
}

func TestMCPToolsCallListUnits(t *testing.T) {
	h := New()
	body := strings.NewReader(`{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_units","arguments":{"category":"length"}},"id":5}`)
	req := httptest.NewRequest("POST", "/mcp", body)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("MCP list_units status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("MCP list_units parse error: %v, body: %s", err, w.Body.String())
	}
	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		t.Fatal("MCP list_units should return result")
	}
	content, ok := result["content"].([]interface{})
	if !ok || len(content) == 0 {
		t.Fatal("MCP list_units should return content array")
	}
	textItem, _ := content[0].(map[string]interface{})
	text, _ := textItem["text"].(string)
	if !strings.Contains(text, "category=length") {
		t.Errorf("MCP list_units text should contain category=length, got: %s", text)
	}
}

func TestMCPMethodNotFound(t *testing.T) {
	h := New()
	body := strings.NewReader(`{"jsonrpc":"2.0","method":"unknown/method","id":6}`)
	req := httptest.NewRequest("POST", "/mcp", body)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("MCP parse error: %v, body: %s", err, w.Body.String())
	}
	if resp["error"] == nil {
		t.Error("MCP unknown method should return error")
	}
}
