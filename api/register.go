package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/khorsl/teacherportal/db/sqlc"
)

type createRegistersRequest struct {
	Teacher  string   `json:"teacher" binding:"required,email"`
	Students []string `json:"students" binding:"required,dive,email"`
}

func (server *Server) createRegisters(ctx *gin.Context) {
	var req createRegistersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.RegisterTxParams{
		Teacher:  req.Teacher,
		Students: req.Students,
	}

	err := server.store.RegisterTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusNoContent)
}
