package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mrExplorist/shorten-url-fiber-redis/database"
	"github.com/mrExplorist/shorten-url-fiber-redis/helpers"
	"github.com/redis/go-redis/v9"
)

// Define the request and response structures
type request struct {
	URL 								string `json:"url"`
	CustomShort 				string `json:"custom_short"`
	Expiry 							time.Duration `json:"expiry"`
} 

type response struct {
	URL   								string  `json:"url"`
	CustomShort  					string   `json:"custom_short"`
	Expiry   							time.Duration  `json:"expiry"`
	XRateRemaining				int           `json:"rate_limit"`
	XRateLimitReset				time.Duration  `json:"rate_limit_reset"`
} 


func ShortenURL(c *fiber.Ctx) error{
	// 1. Parse the request body 
	// 2. Validate the request body
	// 3. Check if the URL is already shortened
	// 4. If not, shorten the URL
	// 5. Return the response

	body := new(request)
	// Parse the request body so golang can understand it into struct
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ "error": true, "message": "Cannot parse JSON", "data": nil })
	}

	

	// ~---*****************************--->

	// Implement rate limiting here 
	// We will use the IP address of the user to implement rate limiting 
	// We will allow 10 requests per 30 minutes

	// ~---*****************************--->



//TODO:  implement the rate limiting here - 10 requests per 30 minutes --------------------**********************>
	r2 := database.NewRedisClient(2) // create a new redis client
	defer r2.Close() // close the redis client
  _ , err := r2.Get(database.Ctx, c.IP()).Result() // IP address of the user is the key in the redis server 
	if err == redis.Nil {// if the key is not found in the redis server 
		// set the key API in the redis server with the value API_Quota  and the expiry of 30 minutes 
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("RATE_LIMIT"), 30 * time.Minute).Err()
	} else {
		
		// if the key is found in the redis server
		// get the value of the key from the redis server
		// convert the value to int
		// if the value is greater than 10, return an error
		// if the value is less than 10, increment the value by 1 and set the key in the redis server with the new value
		// if the value is equal to 10, return an error

		val , _ := r2.Get(database.Ctx, c.IP()).Result() // get the value of the key from the redis server 
		valInt , _ := strconv.Atoi(val) // convert the value to int 
		if valInt <= 0 { 
			limit , _ := r2.TTL(database.Ctx, c.IP()).Result() // get the time to live of the key from the redis server
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{ "error": true, "message": "Rate limit exceeded. Try again in " + strconv.Itoa(int(limit.Seconds())) + " seconds", "data": nil , "X-RateLimit-Remaining": val, "X-RateLimit-Reset": limit / time.Nanosecond / time.Minute })
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ "error": true, "message": "Something went wrong in connecting to the DB", "data": nil })
		}
	}



	


	// Check if input is an actual URL
	// We will use the go validator package to validate the URL
	// govalidator.IsURL() returns true if the input is a valid URL and false if it is not a valid URL 
	
	 if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ "error": true, "message": "Invalid URL", "data": nil })
	}


	// check for domain error 

	if !helpers.IsDomainValid(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{ "error": true, "message": "Invalid URL", "data": nil })
	}


	// enforce http, SSL

	body.URL = helpers.EnforceHTTP(body.URL)
	// check if the URL is already shortened
	// we will use the redis database to check if the URL is already shortened
	// if the URL is already shortened, we will return the shortened URL
	// if the URL is not shortened, we will shorten the URL and return the shortened URL


	// Accept the custom short URL if it is provided by the user 
	// If the custom short URL is not provided by the user, we will generate a random short URL

	var id string 
	if body.CustomShort != "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.NewRedisClient(0) // create a new redis client
	defer r.Close() // close the redis client

	// check if the customURL is already taken 
	// if the customURL is already taken, return an error
	// if the customURL is not taken, shorten the URL and return the shortened URL

	val, _ :=  r.Get(database.Ctx, id).Result() // get the value of id from the redis server


	if val != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ "error": true, "message": "Custom URL already taken", "data": nil })
	}

	// check if the expiry is provided by the user
	// if the expiry is provided by the user, set the key in the redis server with the expiry
	

	if body.Expiry != 0 {
		body.Expiry = 24
	}
 
	// set the key in the redis server with the value of the URL and the expiry
 err = r.Set(database.Ctx, id, body.URL, body.Expiry * time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ "error": true, "message": "Something went wrong in connecting to the Server", "data": nil })
	}


	



	// handling the response 
	// return the shortened URL to the user
	resp := response {


		URL : body.URL,
		CustomShort: " ",
		Expiry: body.Expiry,
		XRateRemaining : 10,
		XRateLimitReset : 30 * time.Minute,

	
	
	}





	r2.Decr(database.Ctx, c.IP()) // decrement the value of the key in the redis server by 1

	val, _ = r2.Get(database.Ctx, c.IP()).Result() // get the value of the key from the redis server
 
	r2.Decr(database.Ctx, c.IP()) // decrement the value of the key in the redis server by 1

	resp.XRateRemaining , _ = strconv.Atoi(val) // convert the value to int and assign it to the XRateRemaining field in the response struct

	ttl , _ := r2.TTL(database.Ctx, c.IP()).Result() // get the time to live of the key from the redis server and assign it to the XRateLimitReset field in the response struct

	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute // convert the time to live to minutes


	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id // assign the shortened URL to the CustomShort field in the response struct

	return c.Status(fiber.StatusOK).JSON(fiber.Map{ "error": false, "message": "URL shortened successfully", "data": resp })


}


