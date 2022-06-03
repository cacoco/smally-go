package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

const initialValue int64 = 10000000
const counterKey string = "smally:counter"
const urlKeyPrefix string = "url-"
const encodingRadix = 32

type (
	postURLRequest struct {
		URL string `json:"url"`
	}

	postURLResponse struct {
		SmallyURL string `json:"smally_url"`
	}

	// for testing
	handler struct {
		rdb *redis.Client
	}
)

func newPostURLResponse(url string) *postURLResponse {
	p := postURLResponse{SmallyURL: url}
	return &p
}

var redisCtx = context.Background()
var redisClient *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

//----------
// Handlers
//----------

// Create - POST /url Handler
func Create(c echo.Context) error {
	h := &handler{redisClient}
	return h.createSmallyURL(c)
}

func (h *handler) createSmallyURL(c echo.Context) error {
	p := new(postURLRequest)
	if err := c.Bind(p); err != nil {
		return err
	}

	fmt.Printf("Creating shortened url for %s\n", p.URL)

	next := h.nextCounter()
	redisURLKey := fmt.Sprintf("%s%s", urlKeyPrefix, fmt.Sprintf("%d", next))
	err := h.rdb.Set(redisCtx, redisURLKey, p.URL, 0).Err()
	if err != nil {
		return err
	}

	id := strconv.FormatInt(next, encodingRadix)
	smallyURL := fmt.Sprintf("%s://%s/%s", c.Scheme(), c.Request().Host, id)

	return c.JSON(http.StatusCreated, newPostURLResponse(smallyURL))
}

// Get - GET /:id Handler
func Get(c echo.Context) error {
	h := &handler{redisClient}
	return h.getSmallyURL(c)
}

func (h *handler) getSmallyURL(c echo.Context) error {
	id := c.Param("id")

	fmt.Printf("Decoding url for shortened id %s\n", id)

	decodedID, err := strconv.ParseInt(id, encodingRadix, 64)
	if err != nil {
		return err
	}
	redisURLKey := fmt.Sprintf("%s%s", urlKeyPrefix, fmt.Sprintf("%d", decodedID))

	url, err := h.rdb.Get(redisCtx, redisURLKey).Result()
	if err != nil {
		fmt.Printf("No matching url or encountered an error looking up matching url. %s\n", err.Error())
		return err
	}
	fmt.Printf("Found matching url: %s for id: %s\n", url, id)
	return c.Redirect(http.StatusMovedPermanently, url)
}

//----------
// Services
//----------

func (h *handler) getShortURL(id string) string {
	decodedID, err := strconv.ParseInt(id, encodingRadix, 64)
	if err != nil {
		panic(err)
	}
	redisURLKey := fmt.Sprintf("%s%s", urlKeyPrefix, fmt.Sprintf("%d", decodedID))

	url, err := h.rdb.Get(redisCtx, redisURLKey).Result()
	if err != nil {
		panic(err)
	} else {
		return url
	}
}

func (h *handler) getCounter() int64 {
	count, err := h.rdb.Get(redisCtx, counterKey).Int64()
	if err == redis.Nil {
		return initialValue
	} else if err != nil {
		panic(err)
	} else {
		return count
	}
}

func (h *handler) incrCounter(val int64) {
	next := val + 1
	err := h.rdb.Set(redisCtx, counterKey, fmt.Sprintf("%d", next), 0).Err()
	if err != nil {
		panic(err)
	}
}

func (h *handler) nextCounter() int64 {
	current := h.getCounter()
	h.incrCounter(current)
	return current
}
