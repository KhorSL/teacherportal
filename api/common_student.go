package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/khorsl/teacherportal/db/sqlc"
	"github.com/khorsl/teacherportal/util"
)

type commonStudentsRequest struct {
	TeacherEmails []string `form:"teacher" binding:"required,dive,email"`
}

type commonStudentsResponse struct {
	Students []string `json:"students"`
}

func newCommonStudentsResponse(emails []string) commonStudentsResponse {
	return commonStudentsResponse{
		Students: emails,
	}
}

func (server *Server) commonStudent(ctx *gin.Context) {
	var req commonStudentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetCommonStudentsEmailParams{
		Email: util.StringsToSingleQuoteCommaSep(req.TeacherEmails),
		Count: int64(len(req.TeacherEmails)),
	}

	emails, err := server.store.GetCommonStudentsEmail(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := newCommonStudentsResponse(emails)
	ctx.JSON(http.StatusOK, response)
}
