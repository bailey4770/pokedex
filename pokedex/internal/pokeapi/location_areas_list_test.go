package pokeapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Helper function to create a client with zero timeout for testing
func testClient() Client {
	return NewClient(0)
}

func TestGetLocationData_Success(t *testing.T) {
	// Mock server returning valid JSON
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"count": 2,
			"next": "next-url",
			"previous": "prev-url",
			"results": [
				{"name": "area1", "url": "url1"},
				{"name": "area2", "url": "url2"}
			]
		}`))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	client := testClient()
	data, err := client.GetLocationData(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(data.Results))
	}
	if data.Results[0].Name != "area1" {
		t.Errorf("expected first result 'area1', got %s", data.Results[0].Name)
	}
	if data.Next == nil || *data.Next != "next-url" {
		t.Errorf("expected next-url, got %v", data.Next)
	}
	if data.Previous == nil || *data.Previous != "prev-url" {
		t.Errorf("expected prev-url, got %v", data.Previous)
	}
}

func TestGetLocationData_HTTPError(t *testing.T) {
	// Mock server returning 500
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`internal server error`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := testClient()
	_, err := client.GetLocationData(server.URL)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if !contains(err.Error(), "response failed with status code 500") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetLocationData_InvalidJSON(t *testing.T) {
	// Mock server returning invalid JSON
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid-json}`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := testClient()
	_, err := client.GetLocationData(server.URL)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if !contains(err.Error(), "unmarshalling response body") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetLocationData_Timeout(t *testing.T) {
	// Mock server that sleeps longer than the client timeout
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create client with short timeout
	client := NewClient(10 * time.Millisecond)
	_, err := client.GetLocationData(server.URL)
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
	if !contains(err.Error(), "Client.Timeout") {
		t.Errorf("unexpected error: %v", err)
	}
}

// Helper function to check if a string is contained in another
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || (len(substr) > 0 && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
