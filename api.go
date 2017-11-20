package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func ApiAttendeesHandler(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "https://kumarathonbkk.bookzy.co.th")
	type Attendee struct {
		ID        int
		OrderID   int
		Firstname string
		Lastname  string
		ItemName  string
	}
	var attendees []Attendee

	key := "api_attendee"
	redisClient := OpenRedis()
	results, err := redisClient.Get(key).Result()
	if err != nil || err == redis.Nil {
		db, _ := OpenDB()
		db.Raw("SELECT id, order_id, firstname, lastname, sku item_name FROM attendees ORDER BY order_id, id").Scan(&attendees)
		c.JSON(http.StatusOK, attendees)

		b, err := json.Marshal(attendees)
		if err != nil {
			log.Println(err)
		}

		err = redisClient.Set(key, b, 10*time.Minute).Err()
		if err != nil {
			log.Println(err)
		}
		return
	}

	json.Unmarshal([]byte(results), &attendees)
	c.JSON(http.StatusOK, attendees)
}
