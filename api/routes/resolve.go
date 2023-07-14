package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mrExplorist/shorten-url-fiber-redis/database"
	"github.com/redis/go-redis/v9"
)

// ResolveURL resolves the short url to the original url
func ResolveURL(c *fiber.Ctx ) error {
	url := c.Params("url") // get the url from the params

	// 1. Create a new redis client 
	// 2. Get the url from the redis server
	// 3. Redirect to the url
	// 4. Return an error if the url is not found


	r := database.NewRedisClient(0) // create a new redis client
	defer r.Close() // close the redis client

	val, err := r.Get(database.Ctx, url).Result() // get the url from the redis server

	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ShortURL not found in the DB",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong in connecting to the DB",
		})
	}

	rInr := database.NewRedisClient(1) // create a new redis client
	defer rInr.Close() // close the redis client


	_ = rInr.Incr(database.Ctx,"counter") // increment the url counter in the redis server 


return c.Redirect(val , 301) // redirect to the url



}