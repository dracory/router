package main_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/rtr"
)

func TestRouterEndpoints(t *testing.T) {
	// Create a new router instance with test routes
	r := setupTestRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET /hello",
			method:         http.MethodGet,
			path:           "/hello",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello, World!",
		},
		{
			name:           "GET /api/status",
			method:         http.MethodGet,
			path:           "/api/status",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status": "ok"}`,
		},
		{
			name:           "GET /api/users",
			method:         http.MethodGet,
			path:           "/api/users",
			expectedStatus: http.StatusOK,
			expectedBody:   "List of users",
		},
		{
			name:           "GET /api/users/123",
			method:         http.MethodGet,
			path:           "/api/users/123",
			expectedStatus: http.StatusOK,
			expectedBody:   "User ID: 123",
		},
		{
			name:           "Non-existent route",
			method:         http.MethodGet,
			path:           "/nonexistent",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
			if !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestRouterMethodNotAllowed(t *testing.T) {
	r := setupTestRouter()

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"POST /hello (not allowed)", http.MethodPost, "/hello"},
		{"PUT /api/status (not allowed)", http.MethodPut, "/api/status"},
		{"DELETE /api/users (not allowed)", http.MethodDelete, "/api/users"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			// The current implementation returns 404 for non-matching methods
			// This could be updated to return 405 Method Not Allowed in the future
			if rr.Code != http.StatusNotFound {
				t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
			}
		})
	}
}

// setupTestRouter creates a router instance with test routes
func setupTestRouter() rtr.RouterInterface {
	r := rtr.NewRouter()

	// Add a simple route
	r.AddRoute(rtr.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))

	// Create an API group
	api := rtr.NewGroup().SetPrefix("/api")

	// Add routes to the API group
	api.AddRoute(rtr.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	}))

	// Create a users group with nested routes
	users := rtr.NewGroup().SetPrefix("/users")

	// Add user routes
	users.AddRoute(rtr.Get("", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("List of users"))
	}))

	// Example of a specific user route (exact match required)
	users.AddRoute(rtr.Get("/123", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User ID: 123"))
	}))

	// Add the users group to the API group
	api.AddGroup(users)

	// Add the API group to the router
	r.AddGroup(api)

	return r
}
