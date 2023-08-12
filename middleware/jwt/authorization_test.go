package jwt

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	validToken string

	r *gin.Engine
)

func testHandler(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "ok",
	})
}

func init() {
	validToken, _ = GenerateToken(1, "test")

	r = gin.Default()
	r.GET("/test_get/", AuthorizeGet(), testHandler)
	r.POST("/test_post/", AuthorizePost(), testHandler)
	r.GET("/test_feed/", AuthorizeFeed(), testHandler)
}

func TestAuthorizeGetWithInvalidToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/test_get/?token="+validToken+"1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthorizeGetWithValidToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/test_get/?token="+validToken, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAuthorizePostWithInvalidToken(t *testing.T) {
	req := httptest.NewRequest("POST", "/test_post/?token="+validToken+"1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthorizePostWithValidToken(t *testing.T) {
	req := httptest.NewRequest("POST", "/test_post/?token="+validToken, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAuthorizeFeedWithEmptyToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/test_feed/", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAuthorizeFeedWithInvalidToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/test_feed/?token="+validToken+"1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthorizeFeedWithValidToken(t *testing.T) {
	req := httptest.NewRequest("GET", "/test_feed/?token="+validToken, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
