package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/igorextrawest/GolangContainers/src/constants"
	"github.com/igorextrawest/GolangContainers/src/models"
	"github.com/nats-io/go-nats"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"strings"
)

func getUsersHandler(ctx *gin.Context) {
	redisConn, ok := ctx.MustGet("redis").(*redis.Pool)
	if !ok {
		ctx.String(http.StatusForbidden, "Can't get redis connection")
		return
	}

	query := ctx.Query("query")
	offset, err := getTimeOffset(query)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	conn := redisConn.Get()
	defer conn.Close()

	to := time.Now().Unix()
	from := to - int64(offset)

	toStr := strconv.FormatInt(to, 10)
	fromStr := strconv.FormatInt(from, 10)

	if values, err := redis.Strings(conn.Do(constants.ZRANGEBYSCORE, constants.NamesKey, fromStr, toStr)); err != nil {
		ctx.String(http.StatusForbidden, "Can't get redis connection")
		return
	} else {
		result := []string{}
		for _, val := range values {
			result = append(result, val)
		}
		ctx.JSON(http.StatusOK, &models.Response{Values: result})
	}
}

func createUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.String(http.StatusForbidden, "Can't parse request: %s", err.Error())
		return
	}

	user.Timestamp = time.Now().Unix()
	log.Printf("New User: %+v", user)

	natsConn, ok := ctx.MustGet("nats").(*nats.EncodedConn)
	if !ok {
		ctx.String(http.StatusForbidden, "Can't get nats connection")
		return
	}

	if err := natsConn.Publish(constants.MessageTopic, &user); err != nil {
		ctx.String(http.StatusForbidden, "Can't publish value")
	} else {
		ctx.String(http.StatusOK, "Successfully added!")
	}
}

func getTimeOffset(query string) (int, error) {
	pattern := `(\d*)(\D)`
	pathMetadata := regexp.MustCompile(pattern)

	matches := pathMetadata.FindStringSubmatch(query)

	if len(matches) < 3 {
		return 0, errors.New(fmt.Sprintf("Query not valid: %s", query))
	}

	if matches[1] == "" || matches[2] == "" {
		return 0, errors.New(fmt.Sprintf("Cant's parse query: %s", query))
	}

	numVal, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Cant's parse numerical value: %s", err.Error()))
	}

	mult := getMultiplier(matches[2])
	if mult == 0 {
		return 0, errors.New("Cant's parse multiplier format")
	}

	result := numVal * mult
	return result, nil
}

func getMultiplier(val string) int {
	var multiplier int

	switch strings.ToLower(val) {
	case "s":
		multiplier = 1
		break
	case "m":
		multiplier = 60
		break
	case "h":
		multiplier = 60 * 60
	default:
		multiplier = 0
	}

	return multiplier
}
