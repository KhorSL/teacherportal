package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/khorsl/teacherportal/db/sqlc"
	"github.com/khorsl/teacherportal/util"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store, config util.Config) *Server {
	server := &Server{store: store}
	server.setupRouter(config)
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}

func errorResponseCustom(message string) gin.H {
	return gin.H{"message": message}
}
