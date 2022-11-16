package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/khorsl/teacherportal/db/sqlc"
)

type createTeacherRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type teacherResponse struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func newTeacherResponse(teacher db.Teacher) teacherResponse {
	return teacherResponse{
		ID:        teacher.ID,
		FullName:  teacher.FullName,
		Email:     teacher.Email,
		IsActive:  teacher.IsActive,
		CreatedAt: teacher.CreatedAt,
		UpdatedAt: teacher.UpdatedAt,
		DeletedAt: teacher.DeletedAt,
	}
}

func (server *Server) createTeacher(ctx *gin.Context) {
	var req createTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateTeacherParams{
		FullName: req.FullName,
		Email:    req.Email,
	}

	err := server.store.CreateTeacher(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	teacher, err := server.store.GetTeacherByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newTeacherResponse(teacher)
	ctx.JSON(http.StatusOK, response)
}
