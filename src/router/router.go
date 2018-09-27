package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/nats-io/go-nats"
)

const (
	v1      = "/v1"
	userURL = "/user"
)

type server struct {
	Redis *redis.Pool       `inject:""`
	Nats  *nats.EncodedConn `inject:""`
}

func NewServer() *server {
	return new(server)
}

func (s *server) Start() {
	router := gin.Default()
	router.Use(s.restMiddleware())
	v1 := router.Group(v1)
	{
		v1.POST(userURL, createUser)
		v1.GET(userURL, getUsersHandler)
	}

	router.Run(":8080")
}

func (s *server) restMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("nats", s.Nats)
		c.Set("redis", s.Redis)
		c.Next()
	}
}
