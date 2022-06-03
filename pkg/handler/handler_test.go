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

const longUrl = "https://www.nytimes.com/2022/06/01/well/eat/coffee-study-lower-dying-risk.html"
const smallUrl = "http://example.com/9h5k0"

var (
	postUrlRequestJSON  = `{"url":"https://www.nytimes.com/2022/06/01/well/eat/coffee-study-lower-dying-risk.html"}\n`
	postUrlResponseJSON = fmt.Sprintf(`{"smally_url":"%s"}
`, smallUrl)
)

func TestCreateSmallyUrl(t *testing.T) {
	db, mock := redismock.NewClientMock()

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/url", strings.NewReader(postUrlRequestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//.Get(redisCtx, counterKey).Int64()
	mock.ExpectGet(counterKey).RedisNil()
	next := 10000001
	//.Set(redisCtx, counterKey, next, 0).Err()
	mock.ExpectSet(counterKey, fmt.Sprintf("%d", next), 0).SetVal(fmt.Sprintf("%d", next))
	//.Set(redisCtx, redisUrlKey, p.Url, 0).Err()
	redisUrlKey := fmt.Sprintf("%s%s", urlKeyPrefix, fmt.Sprintf("%d", initialValue))
	mock.ExpectSet(redisUrlKey, longUrl, 0).SetVal(longUrl)

	h := &handler{db}
	// Assertions
	if assert.NoError(t, h.CreateSmallyUrl(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, postUrlResponseJSON, rec.Body.String())
	}
}

func TestGetSmallyUrl(t *testing.T) {
	db, mock := redismock.NewClientMock()

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/:id")
	c.SetParamNames("id")
	c.SetParamValues("9h5k0")

	redisUrlKey := fmt.Sprintf("%s%s", urlKeyPrefix, fmt.Sprintf("%d", initialValue))
	//.Get(redisCtx, redisUrlKey).Result()
	mock.ExpectGet(redisUrlKey).SetVal(longUrl)

	h := &handler{db}
	// Assertions
	if assert.NoError(t, h.GetSmallyUrl(c)) {
		assert.Equal(t, http.StatusMovedPermanently, rec.Code)
		assert.Equal(t, longUrl, rec.HeaderMap.Get("Location"))
	}
}
