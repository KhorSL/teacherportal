package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type suspendStudentRequest struct {
	StudentEmail string `json:"student" binding:"required,email"`
}

func (server *Server) suspendStudent(ctx *gin.Context) {
	var req suspendStudentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := server.store.GetStudentByEmail(ctx, req.StudentEmail)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			ctx.JSON(http.StatusNotFound, errorResponseCustom("student does not exist."))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.GetNotSuspendedStudentByEmail(ctx, req.StudentEmail)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			ctx.JSON(http.StatusBadRequest, errorResponseCustom("student already suspended."))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.SuspendStudentByEmail(ctx, req.StudentEmail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
