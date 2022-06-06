package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const longURL = "https://www.nytimes.com/2022/06/01/well/eat/coffee-study-lower-dying-risk.html"
const smallURL = "http://example.com/9h5k0"
const postURLRequestJSON = `{"url":"https://www.nytimes.com/2022/06/01/well/eat/coffee-study-lower-dying-risk.html"}\n`

var (
	postURLResponseJSON = fmt.Sprintf(`{"smally_url":"%s"}
`, smallURL)
)

func TestCreateSmallyURL(t *testing.T) {
	db, mock := redismock.NewClientMock()

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/url", strings.NewReader(postURLRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//.Get(redisCtx, counterKey).Int64()
	mock.ExpectGet(counterKey).RedisNil()
	next := 10000001
	//.Set(redisCtx, counterKey, next, 0).Err()
	mock.ExpectSet(counterKey, fmt.Sprintf("%d", next), 0).SetVal(fmt.Sprintf("%d", next))
	//.Set(redisCtx, redisURLKey, p.Url, 0).Err()
	redisURLKey := fmt.Sprintf("%s%s", urlKeyPrefix, fmt.Sprintf("%d", initialValue))
	mock.ExpectSet(redisURLKey, longURL, 0).SetVal(longURL)

	h := &handler{db}
	// Assertions
	if assert.NoError(t, h.createSmallyURL(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, postURLResponseJSON, rec.Body.String())
	}
}

func TestGetSmallyURL(t *testing.T) {
	db, mock := redismock.NewClientMock()

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("9h5k0")

	redisURLKey := fmt.Sprintf("%s%s", urlKeyPrefix, fmt.Sprintf("%d", initialValue))
	//.Get(redisCtx, redisUrlKey).Result()
	mock.ExpectGet(redisURLKey).SetVal(longURL)

	h := &handler{db}
	// Assertions
	if assert.NoError(t, h.getSmallyURL(c)) {
		assert.Equal(t, http.StatusMovedPermanently, rec.Code)
		assert.Equal(t, longURL, rec.HeaderMap.Get("Location"))
	}
}
