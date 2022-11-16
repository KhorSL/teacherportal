package api

import (
	"github.com/gin-gonic/gin"
	"github.com/khorsl/teacherportal/util"
)

func routePrefix(route string, prefix string) string {
	return prefix + route
}

func (server *Server) setupRouter(config util.Config) {
	router := gin.Default()
	prefix := config.APIPrefix

	router.POST(routePrefix("/teachers", prefix), server.createTeacher)

	router.POST(routePrefix("/students", prefix), server.createStudent)

	router.POST(routePrefix("/register", prefix), server.createRegisters)

	router.POST(routePrefix("/suspend", prefix), server.suspendStudent)

	router.GET(routePrefix("/commonstudents", prefix), server.commonStudent)

	router.POST(routePrefix("/retrievefornotifications", prefix), server.retreiveStudentEmailsForNotifcations)

	server.router = router
}
