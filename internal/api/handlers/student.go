package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/75asu/sre-bootcamp-one2n/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StudentHandler struct {
	db  *sql.DB
	log *zap.Logger
}

func NewStudentHandler(db *sql.DB, log *zap.Logger) *StudentHandler {
	return &StudentHandler{db: db, log: log}
}

func (h *StudentHandler) Create(c *gin.Context) {
	var req models.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid create request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var student models.Student
	err := h.db.QueryRow(
		`INSERT INTO students (name, email, age, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, name, email, age, created_at, updated_at`,
		req.Name, req.Email, req.Age, time.Now(), time.Now(),
	).Scan(&student.ID, &student.Name, &student.Email, &student.Age, &student.CreatedAt, &student.UpdatedAt)
	if err != nil {
		h.log.Error("failed to create student", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("student created", zap.Int("id", student.ID))
	c.JSON(http.StatusCreated, student)
}

func (h *StudentHandler) GetAll(c *gin.Context) {
	rows, err := h.db.Query(
		`SELECT id, name, email, age, created_at, updated_at FROM students`,
	)
	if err != nil {
		h.log.Error("failed to fetch students", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	students := []models.Student{}
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.Age, &s.CreatedAt, &s.UpdatedAt); err != nil {
			h.log.Error("failed to scan student row", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		students = append(students, s)
	}

	c.JSON(http.StatusOK, students)
}

func (h *StudentHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var student models.Student
	err = h.db.QueryRow(
		`SELECT id, name, email, age, created_at, updated_at FROM students WHERE id = $1`, id,
	).Scan(&student.ID, &student.Name, &student.Email, &student.Age, &student.CreatedAt, &student.UpdatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}
	if err != nil {
		h.log.Error("failed to fetch student", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid update request", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var student models.Student
	err = h.db.QueryRow(
		`UPDATE students SET name=$1, email=$2, age=$3, updated_at=$4
		 WHERE id=$5
		 RETURNING id, name, email, age, created_at, updated_at`,
		req.Name, req.Email, req.Age, time.Now(), id,
	).Scan(&student.ID, &student.Name, &student.Email, &student.Age, &student.CreatedAt, &student.UpdatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}
	if err != nil {
		h.log.Error("failed to update student", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("student updated", zap.Int("id", student.ID))
	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result, err := h.db.Exec(`DELETE FROM students WHERE id = $1`, id)
	if err != nil {
		h.log.Error("failed to delete student", zap.Int("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	h.log.Info("student deleted", zap.Int("id", id))
	c.JSON(http.StatusNoContent, nil)
}
