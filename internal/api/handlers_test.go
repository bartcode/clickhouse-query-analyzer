package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nimbleflux/clickhouse-query-analyzer/internal/clickhouse"
)

func TestWriteJSON_Nil(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusOK, nil)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "null" {
		t.Errorf("expected 'null', got %q", w.Body.String())
	}
}

func TestWriteJSON_NilSlice(t *testing.T) {
	var s []string
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusOK, s)

	if w.Body.String() != "[]" {
		t.Errorf("expected '[]', got %q", w.Body.String())
	}
}

func TestWriteJSON_EmptySlice(t *testing.T) {
	s := []string{}
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusOK, s)

	var result []string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestWriteJSON_ValidData(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, http.StatusOK, map[string]string{"key": "value"})

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json content type")
	}

	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("expected key=value, got %v", result)
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, http.StatusBadRequest, "something went wrong")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var result map[string]string
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if result["error"] != "something went wrong" {
		t.Errorf("expected error message, got %v", result)
	}
}

func TestClientFromRequest_MissingURL(t *testing.T) {
	api := New(clickhouse.NewPool())
	req := httptest.NewRequest("GET", "/", nil)

	_, err := api.clientFromRequest(req)
	if err == nil {
		t.Error("expected error for missing URL")
	}
	if !strings.Contains(err.Error(), "not configured") {
		t.Errorf("expected 'not configured' error, got %v", err)
	}
}

func TestClientFromRequest_Defaults(t *testing.T) {
	api := New(clickhouse.NewPool())
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-CH-URL", "clickhouse://nonexistent:9000")

	_, err := api.clientFromRequest(req)
	if err == nil {
		t.Error("expected connection error for nonexistent host")
	}
}
