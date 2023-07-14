package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background() // create a new context for the redis client 

// 1. Create a new redis client 
// 2. Ping the redis server to check if it is running
// 3. Create a new redis client with
//    a. the host
//    b. the port
//    c. the password
//    d. the database
// 4. Ping the redis server to check if it is running
// 5. Return the redis client



// *redis.Client is a pointer to a redis client 
// redis.Client is a struct that contains the redis client 



func NewRedisClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr : os.Getenv("DB_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"), // password set in the redis server
		DB:       dbNo, // use default DB
	})
		return rdb
}