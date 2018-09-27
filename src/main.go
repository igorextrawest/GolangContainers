package main

import (
	"github.com/facebookgo/inject"
	"github.com/gomodule/redigo/redis"
	"github.com/igorextrawest/GolangContainers/src/consumer"
	"github.com/igorextrawest/GolangContainers/src/router"
	"github.com/nats-io/go-nats"
	"log"
	"time"
)

func main() {
	server := router.NewServer()
	consumer := consumer.NewConsumer()

	err := inject.Populate(
		ConnectNats(),
		ConnectRedis(),
		consumer,
		server,
	)
	if err != nil {
		log.Fatalf("Can't inject values %s", err.Error())
	}

	consumer.Start()
	server.Start()
}

func ConnectNats() *nats.EncodedConn {
	if nc, err := nats.Connect("nats"); err != nil {
		log.Fatalf(err.Error())
		return nil
	} else {
		if c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER); err != nil {
			log.Fatalf(err.Error())
			return nil
		} else {
			return c
		}
	}
}

func ConnectRedis() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: time.Minute,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "redis:6379") },
	}
}
