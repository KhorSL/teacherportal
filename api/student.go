package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/khorsl/teacherportal/db/sqlc"
)

type createStudentRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type studentResponse struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func newStudentResponse(student db.Student) studentResponse {
	return studentResponse{
		ID:        student.ID,
		FullName:  student.FullName,
		Email:     student.Email,
		IsActive:  student.IsActive,
		CreatedAt: student.CreatedAt,
		UpdatedAt: student.UpdatedAt,
		DeletedAt: student.DeletedAt,
	}
}

func (server *Server) createStudent(ctx *gin.Context) {
	var req createStudentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateStudentParams{
		FullName: req.FullName,
		Email:    req.Email,
	}

	err := server.store.CreateStudent(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	student, err := server.store.GetStudentByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newStudentResponse(student)
	ctx.JSON(http.StatusOK, response)
}
