package parse

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestParseQueryParamsWithInvalidParam(t *testing.T) {
	req := httptest.NewRequest("GET", "/test/?p=abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res := getResponse(w)

	assert.Equal(t, int32(0), res.StatusCode)
}

func TestParseQueryParamsWithInvalidIntParam(t *testing.T) {
	req := httptest.NewRequest("GET", "/test/?user_id=abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res := getResponse(w)

	assert.NotEqual(t, int32(0), res.StatusCode)
	assert.Equal(t, "invalid syntax", res.StatusMsg)
}

func TestParseQueryParams(t *testing.T) {
	req := httptest.NewRequest("GET", "/test/?user_id=1&comment_text=abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res := getResponse(w)

	assert.Equal(t, int32(0), res.StatusCode)
}
