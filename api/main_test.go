package api

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/khorsl/teacherportal/db/sqlc"
	"github.com/khorsl/teacherportal/util"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	server := NewServer(store, config)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
