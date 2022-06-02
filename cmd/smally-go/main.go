package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/go-redis/redis/v8"
)

const initialValue int64 = 10000000
const counterKey string = "smally:counter"
const urlKeyPrefix string = "url-"
const encodingRadix = 32

type (
	postUrlRequest struct {
		Url string `json:"url"`
	}

	postUrlResponse struct {
		SmallyUrl string `json:"smallyUrl"`
	}
)

var redisCtx = context.Background()
var redisClient *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func newPostUrlResponse(url string) *postUrlResponse {
	p := postUrlResponse{SmallyUrl: url}
	return &p
}

//----------
// Handlers
//----------

func createSmallyUrl(c echo.Context) error {
	p := new(postUrlRequest)
	if err := c.Bind(p); err != nil {
		return err
	}

	fmt.Printf("Creating shortened url for %s\n", p.Url)

	next := nextCounter()
	redisUrlKey := fmt.Sprintf("%s%s", urlKeyPrefix, fmt.Sprintf("%d", next))
	err := redisClient.Set(redisCtx, redisUrlKey, p.Url, 0).Err()
	if err != nil {
		panic(err)
	}

	id := strconv.FormatInt(next, encodingRadix)
	smallyUrl := fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, id)

	return c.JSON(http.StatusCreated, newPostUrlResponse(smallyUrl))
}

func getSmallyUrl(c echo.Context) error {
	id := c.Param("id")

	fmt.Printf("Decoding url for shortened id %s\n", id)

	decodedId, err := strconv.ParseInt(id, encodingRadix, 64)
	if err != nil {
		panic(err)
	}
	redisUrlKey := fmt.Sprintf("%s%s", urlKeyPrefix, fmt.Sprintf("%d", decodedId))

	url, err := redisClient.Get(redisCtx, redisUrlKey).Result()
	if err != nil {
		fmt.Printf("No matching url or encountered an error looking up matching url. %s", err.Error())
		panic(err)
	} else {
		fmt.Printf("Found matching url: %s for id: %s", url, id)
		return c.Redirect(http.StatusMovedPermanently, url)
	}
}

//----------
// Services
//----------

func getCounter() int64 {
	count, err := redisClient.Get(redisCtx, counterKey).Int64()
	if err == redis.Nil {
		return initialValue
	} else if err != nil {
		panic(err)
	} else {
		return count
	}
}

func incrCounter(val int64) {
	next := val + 1
	err := redisClient.Set(redisCtx, counterKey, next, 0).Err()
	if err != nil {
		panic(err)
	}
}

func nextCounter() int64 {
	current := getCounter()
	incrCounter(current)
	return current
}

func getShortUrl(id string) string {
	decodedId, err := strconv.ParseInt(id, encodingRadix, 64)
	if err != nil {
		panic(err)
	}
	redisUrlKey := fmt.Sprintf("%s%s", urlKeyPrefix, fmt.Sprintf("%d", decodedId))

	url, err := redisClient.Get(redisCtx, redisUrlKey).Result()
	if err != nil {
		panic(err)
	} else {
		return url
	}
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/url", createSmallyUrl)
	e.GET("/:id", getSmallyUrl)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
