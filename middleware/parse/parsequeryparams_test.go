package parse

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ghjattu/tiny-tiktok/bloomfilter"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	r *gin.Engine
)

func testHandler(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "ok",
	})
}

func getResponse(w *httptest.ResponseRecorder) *Response {
	res := &Response{}
	bytes, _ := io.ReadAll(w.Result().Body)
	json.Unmarshal(bytes, res)

	return res
}

func init() {
	r = gin.Default()
	r.Use(ParseQueryParams())
	r.GET("/test/", testHandler)
}

func TestParseQueryParams(t *testing.T) {
	t.Run("invalid param", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/?p=abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := getResponse(w)

		assert.Equal(t, int32(0), res.StatusCode)
	})

	t.Run("invalid int param", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/?user_id=abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := getResponse(w)

		assert.NotEqual(t, int32(0), res.StatusCode)
		assert.Equal(t, "invalid syntax", res.StatusMsg)
	})

	t.Run("less than 0 int param", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/?user_id=-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := getResponse(w)

		assert.NotEqual(t, int32(0), res.StatusCode)
	})

	t.Run("key does not belong to any bloom filter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/?action_type=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := getResponse(w)

		assert.Equal(t, int32(0), res.StatusCode)
	})

	t.Run("key belongs to user bloom filter and does not exist in filter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/?user_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := getResponse(w)

		assert.Equal(t, int32(1), res.StatusCode)
	})

	t.Run("key belongs to user bloom filter and exists in filter", func(t *testing.T) {
		bloomfilter.ClearAll()
		bloomfilter.Add(bloomfilter.UserBloomFilter, 1)

		req := httptest.NewRequest("GET", "/test/?user_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := getResponse(w)

		assert.Equal(t, int32(0), res.StatusCode)
	})

	t.Run("key belongs to video bloom filter and exists in filter", func(t *testing.T) {
		bloomfilter.ClearAll()
		bloomfilter.Add(bloomfilter.VideoBloomFilter, 1)

		req := httptest.NewRequest("GET", "/test/?video_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := getResponse(w)

		assert.Equal(t, int32(0), res.StatusCode)
	})

	t.Run("key belongs to comment bloom filter and exists in filter", func(t *testing.T) {
		bloomfilter.ClearAll()
		bloomfilter.Add(bloomfilter.CommentBloomFilter, 1)

		req := httptest.NewRequest("GET", "/test/?comment_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := getResponse(w)

		assert.Equal(t, int32(0), res.StatusCode)
	})

	t.Run("parse int param successfully", func(t *testing.T) {
		bloomfilter.ClearAll()
		bloomfilter.Add(bloomfilter.UserBloomFilter, 1)

		req := httptest.NewRequest("GET", "/test/?user_id=1&comment_text=abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		res := getResponse(w)

		assert.Equal(t, int32(0), res.StatusCode)
	})
}
