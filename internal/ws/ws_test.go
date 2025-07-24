package ws

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleWS(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		headers        map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET request without WebSocket headers should return error",
			method:         "GET",
			headers:        map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad Request",
		},
		{
			name:   "GET request with WebSocket headers should upgrade",
			method: "GET",
			headers: map[string]string{
				"Connection":            "Upgrade",
				"Upgrade":               "websocket",
				"Sec-WebSocket-Version": "13",
				"Sec-WebSocket-Key":     "dGhlIHNhbXBsZSBub25jZQ==",
			},
			expectedStatus: http.StatusSwitchingProtocols,
		},
		{
			name:           "POST request should return method not allowed",
			method:         "POST",
			headers:        map[string]string{},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method Not Allowed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := httptest.NewRequest(tc.method, "/ws", nil)
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}
			w := httptest.NewRecorder()

			// Act
			HandleWS(w, req)

			// Assert
			if w.Code != tc.expectedStatus {
				t.Errorf("HandleWS() status = %v, want %v", w.Code, tc.expectedStatus)
			}

			if tc.expectedBody != "" {
				body := strings.TrimSpace(w.Body.String())
				if !strings.Contains(body, tc.expectedBody) {
					t.Errorf("HandleWS() body = %v, want to contain %v", body, tc.expectedBody)
				}
			}
		})
	}
}

// TDD: This test should fail because we haven't implemented connection management yet
func TestConnectionManager(t *testing.T) {
	// Test that we can track active connections
	manager := NewConnectionManager()

	if manager == nil {
		t.Error("NewConnectionManager() should return non-nil manager")
	}

	// Test initial state
	count := manager.GetConnectionCount()
	if count != 0 {
		t.Errorf("GetConnectionCount() = %v, want 0 for new manager", count)
	}
}

// TDD: This test should fail because we haven't implemented message broadcasting yet
func TestBroadcastMessage(t *testing.T) {
	manager := NewConnectionManager()

	// Test broadcasting to empty connection pool
	err := manager.Broadcast([]byte("test message"))
	if err != nil {
		t.Errorf("Broadcast() to empty pool should not return error, got: %v", err)
	}
}

// TDD: Test connection management interface
func TestConnectionManagerInterface(t *testing.T) {
	manager := NewConnectionManager()

	// Test adding a connection
	connID := "test-conn-1"
	err := manager.AddConnection(connID, nil)
	if err != nil {
		t.Errorf("AddConnection() error = %v, want nil", err)
	}

	// Test connection count
	count := manager.GetConnectionCount()
	if count != 1 {
		t.Errorf("GetConnectionCount() = %v, want 1", count)
	}

	// Test removing connection
	err = manager.RemoveConnection(connID)
	if err != nil {
		t.Errorf("RemoveConnection() error = %v, want nil", err)
	}

	// Test connection count after removal
	count = manager.GetConnectionCount()
	if count != 0 {
		t.Errorf("GetConnectionCount() after removal = %v, want 0", count)
	}
}

// TDD: Test for production-safe error handling
func TestGetConnectionsSafe(t *testing.T) {
	manager := NewConnectionManager()

	// Test getting connections safely
	connections := manager.GetConnections()
	if connections == nil {
		t.Error("GetConnections() should return non-nil slice")
	}
	if len(connections) != 0 {
		t.Errorf("GetConnections() length = %v, want 0 for empty manager", len(connections))
	}
}
