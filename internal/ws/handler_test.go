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
				"Sec-WebSocket-Key":     "test-key",
			},
			expectedStatus: http.StatusSwitchingProtocols,
			expectedBody:   "",
		},
		{
			name:           "POST request should return Method Not Allowed",
			method:         "POST",
			headers:        map[string]string{},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Method Not Allowed",
		},
		{
			name:   "Incomplete WebSocket headers should return Bad Request",
			method: "GET",
			headers: map[string]string{
				"Connection": "Upgrade",
				"Upgrade":    "websocket",
				// Missing Sec-WebSocket-Version and Sec-WebSocket-Key
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad Request",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest(tc.method, "/ws", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Add headers
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler
			HandleWS(w, req)

			// Check status code
			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			// Check response body if expected
			if tc.expectedBody != "" {
				body := strings.TrimSpace(w.Body.String())
				if !strings.Contains(body, tc.expectedBody) {
					t.Errorf("Expected body to contain '%s', got '%s'", tc.expectedBody, body)
				}
			}
		})
	}
}
