package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/pz3-http/internal/storage"
)

func TestTaskHandlers(t *testing.T) {
	store := storage.NewMemoryStore()
	handlers := NewHandlers(store)

	task := store.Create("Test task")

	t.Run("GET /tasks/{id} - existing task", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasks/1", nil)
		rr := httptest.NewRecorder()

		req.URL.Path = "/tasks/1"
		handlers.GetTask(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var responseTask storage.Task
		if err := json.Unmarshal(rr.Body.Bytes(), &responseTask); err != nil {
			t.Fatalf("could not parse response: %v", err)
		}

		if responseTask.ID != task.ID {
			t.Errorf("expected task ID %d, got %d", task.ID, responseTask.ID)
		}
	})

	t.Run("GET /tasks/{id} - non-existing task", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasks/999", nil)
		rr := httptest.NewRecorder()

		req.URL.Path = "/tasks/999"
		handlers.GetTask(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})

	t.Run("GET /tasks/{id} - invalid ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasks/abc", nil)
		rr := httptest.NewRecorder()

		req.URL.Path = "/tasks/abc"
		handlers.GetTask(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("PATCH /tasks/{id} - update task", func(t *testing.T) {
		updateData := map[string]bool{"done": true}
		body, _ := json.Marshal(updateData)

		req := httptest.NewRequest("PATCH", "/tasks/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		req.URL.Path = "/tasks/1"
		handlers.UpdateTask(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var responseTask storage.Task
		if err := json.Unmarshal(rr.Body.Bytes(), &responseTask); err != nil {
			t.Fatalf("could not parse response: %v", err)
		}

		if responseTask.Done != true {
			t.Errorf("expected task Done true, got %v", responseTask.Done)
		}
	})

	t.Run("PATCH /tasks/{id} - missing done field", func(t *testing.T) {
		updateData := map[string]string{"other": "field"}
		body, _ := json.Marshal(updateData)

		req := httptest.NewRequest("PATCH", "/tasks/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		req.URL.Path = "/tasks/1"
		handlers.UpdateTask(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("DELETE /tasks/{id} - existing task", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/tasks/1", nil)
		rr := httptest.NewRecorder()

		req.URL.Path = "/tasks/1"
		handlers.DeleteTask(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
		}

		_, err := store.Get(1)
		if err == nil {
			t.Error("expected task to be deleted, but it still exists")
		}
	})

	t.Run("DELETE /tasks/{id} - non-existing task", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/tasks/999", nil)
		rr := httptest.NewRecorder()

		req.URL.Path = "/tasks/999"
		handlers.DeleteTask(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})
}

func TestCORSMiddleware(t *testing.T) {
	store := storage.NewMemoryStore()
	handlers := NewHandlers(store)

	t.Run("CORS headers present", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/tasks", nil)
		rr := httptest.NewRecorder()

		handler := CORS(http.HandlerFunc(handlers.ListTasks))
		handler.ServeHTTP(rr, req)

		if origin := rr.Header().Get("Access-Control-Allow-Origin"); origin != "GET, POST" {
			t.Errorf("expected CORS header Access-Control-Allow-Origin: GET, POST, got %s", origin)
		}
	})

	t.Run("OPTIONS preflight request", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/tasks", nil)
		rr := httptest.NewRecorder()

		handler := CORS(http.HandlerFunc(handlers.ListTasks))
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code for OPTIONS: got %v want %v", status, http.StatusOK)
		}

		if methods := rr.Header().Get("Access-Control-Allow-Methods"); methods == "" {
			t.Error("expected CORS header Access-Control-Allow-Methods")
		}
	})
}
