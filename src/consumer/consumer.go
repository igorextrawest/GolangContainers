package consumer

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/igorextrawest/GolangContainers/src/constants"
	"github.com/igorextrawest/GolangContainers/src/models"
	"github.com/nats-io/go-nats"
	"log"
	"strconv"
	"time"
)

type consumer struct {
	Redis *redis.Pool       `inject:""`
	Nats  *nats.EncodedConn `inject:""`
}

func NewConsumer() *consumer {
	return new(consumer)
}

func (c *consumer) Start() {
	c.StartSubscriber()
	c.StartCleaner()
}

func (c *consumer) StartSubscriber() {
	_, err := c.Nats.Subscribe(constants.MessageTopic, func(user *models.User) {
		if user != nil {
			c.SaveToRedis(fmt.Sprintf("%s %s", user.FirstName, user.LastName))
			log.Printf("Received a user: %+v\n", user)
		}
	})
	if err != nil {
		log.Printf("Can't subscribe to %s", constants.MessageTopic)
	}

}

func (c *consumer) SaveToRedis(fullName string) {
	connection := c.Redis.Get()
	defer connection.Close()

	t := time.Now().Unix()
	timeStr := strconv.FormatInt(t, 10)
	if _, err := connection.Do(constants.ZADD, constants.NamesKey, timeStr, fullName); err != nil {
		log.Printf("Can't add value to set of names: %s", err.Error())
		return
	} else {
		log.Printf("%s successfully added to the redis", fullName)
	}
}

func (c *consumer) StartCleaner() {
	go func() {
		for {
			time.Sleep(time.Minute)
			c.checkOutdatedNames()
		}
	}()
}

func (c *consumer) checkOutdatedNames() error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Error on cleaning usernames: %+v", err)
		}
	}()

	connection := c.Redis.Get()
	defer connection.Close()

	to := time.Now().Unix() - constants.ONE_HOUR
	from := int64(0)

	toStr := strconv.FormatInt(to, 10)
	fromStr := strconv.FormatInt(from, 10)

	if keys, err := redis.Strings(connection.Do(constants.ZRANGEBYSCORE, constants.NamesKey, fromStr, toStr)); err != nil {
		return err
	} else {
		for _, key := range keys {
			if _, err := connection.Do(constants.ZREM, constants.NamesKey, key); err != nil {
				log.Printf("Can't remove key %s", key)
			} else {
				log.Printf("Removed %s", key)
			}
		}
	}

	return nil
}
