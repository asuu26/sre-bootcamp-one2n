package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupRouter(h *StudentHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/students", h.Create)
		v1.GET("/students", h.GetAll)
		v1.GET("/students/:id", h.GetByID)
		v1.PUT("/students/:id", h.Update)
		v1.DELETE("/students/:id", h.Delete)
	}
	return r
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "age", "created_at", "updated_at"}).
		AddRow(1, "John Doe", "john@example.com", 22, now, now)

	mock.ExpectQuery(`INSERT INTO students`).
		WithArgs("John Doe", "john@example.com", 22, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(rows)

	h := NewStudentHandler(db, zap.NewNop())
	r := setupRouter(h)

	body, _ := json.Marshal(map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   22,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreate_ValidationError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	h := NewStudentHandler(db, zap.NewNop())
	r := setupRouter(h)

	body, _ := json.Marshal(map[string]interface{}{
		"name": "John Doe",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "age", "created_at", "updated_at"}).
		AddRow(1, "John Doe", "john@example.com", 22, now, now).
		AddRow(2, "Jane Doe", "jane@example.com", 23, now, now)

	mock.ExpectQuery(`SELECT id, name, email, age, created_at, updated_at FROM students`).
		WillReturnRows(rows)

	h := NewStudentHandler(db, zap.NewNop())
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "age", "created_at", "updated_at"}).
		AddRow(1, "John Doe", "john@example.com", 22, now, now)

	mock.ExpectQuery(`SELECT id, name, email, age, created_at, updated_at FROM students WHERE id`).
		WithArgs(1).
		WillReturnRows(rows)

	h := NewStudentHandler(db, zap.NewNop())
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	mock.ExpectQuery(`SELECT id, name, email, age, created_at, updated_at FROM students WHERE id`).
		WithArgs(99).
		WillReturnRows(sqlmock.NewRows([]string{}))

	h := NewStudentHandler(db, zap.NewNop())
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/99", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "age", "created_at", "updated_at"}).
		AddRow(1, "Jane Doe", "jane@example.com", 23, now, now)

	mock.ExpectQuery(`UPDATE students`).
		WithArgs("Jane Doe", "jane@example.com", 23, sqlmock.AnyArg(), 1).
		WillReturnRows(rows)

	h := NewStudentHandler(db, zap.NewNop())
	r := setupRouter(h)

	body, _ := json.Marshal(map[string]interface{}{
		"name":  "Jane Doe",
		"email": "jane@example.com",
		"age":   23,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/students/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 got %d: %s", w.Code, w.Body.String())
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	mock.ExpectExec(`DELETE FROM students`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	h := NewStudentHandler(db, zap.NewNop())
	r := setupRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/students/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204 got %d: %s", w.Code, w.Body.String())
	}
}
