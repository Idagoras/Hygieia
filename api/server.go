package api

import (
	"Hygieia/config"
	"Hygieia/database"
	"Hygieia/mail"
	"Hygieia/oss"
	"Hygieia/pm"
	"Hygieia/sms"
	"Hygieia/token"
	"Hygieia/worker"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	router      *gin.Engine
	rdb         *redis.Client
	store       database.Store
	ossService  oss.Service
	enforcer    pm.Enforcer
	sms         sms.Sender
	mailSender  mail.Sender
	tokenMaker  token.Maker
	config      *config.Config
	distributor worker.TaskDistributor
}

func NewServer(rdb *redis.Client, store database.Store) (*Server, error) {
	server := &Server{
		rdb:   rdb,
		store: store,
	}
	router := gin.Default()
	server.router = router
	server.initRouter()
	return server, nil
}

func (server *Server) initRouter() {
	r := server.router
	r.POST("/eegSession/begin", server.BeginEEGSession)

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
