// Description:
// This testutils.go file contains some functions and variables
// that are used in test files in the controllers package.

package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/Ghjattu/tiny-tiktok/middleware/jwt"
	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	serverIP   string
	serverPort string

	r *gin.Engine
)

// init() retrieves the environment variables, initializes the gin engine.
func init() {
	godotenv.Load("../.env")
	serverIP = os.Getenv("SERVER_IP")
	serverPort = os.Getenv("SERVER_PORT")

	r = gin.Default()
	r.GET("/douyin/feed/", jwt.AuthorizeFeed(), Feed)
	r.POST("/douyin/user/register/", Register)
	r.POST("/douyin/user/login/", Login)
	r.GET("/douyin/user/", jwt.AuthorizeGet(), GetUserByUserIDAndToken)
	r.POST("/douyin/publish/action/", jwt.AuthorizePost(), PublishNewVideo)
	r.GET("/douyin/publish/list/", jwt.AuthorizeGet(), GetPublishListByAuthorID)

	r.POST("/douyin/favorite/action/", jwt.AuthorizePost(), FavoriteAction)
	r.GET("/douyin/favorite/list/", jwt.AuthorizeGet(), GetFavoriteListByUserID)
	r.POST("/douyin/comment/action/", jwt.AuthorizePost(), CommentAction)
	r.GET("/douyin/comment/list/", jwt.AuthorizeGet(), CommentList)

	r.POST("/douyin/relation/action/", jwt.AuthorizePost(), FollowAction)
}

// selectResponseType selects the response type according to the request path.
//
//	@param req *http.Request
//	@return interface{}
func selectResponseType(req *http.Request) interface{} {
	path, _ := strings.CutPrefix(req.URL.Path, "/douyin")

	switch path {
	case "/feed/":
		return &FeedResponse{}
	case "/user/register/":
		return &RegisterResponse{}
	case "/user/login/":
		return &LoginResponse{}
	case "/user/":
		return &UserResponse{}
	case "/publish/list/":
		return &PublishListResponse{}
	case "/favorite/list/":
		return &FavoriteListResponse{}
	case "/comment/action/":
		return &CommentActionResponse{}
	case "/comment/list/":
		return &CommentListResponse{}
	default:
		return &Response{}
	}
}

// registerTestUser registers a new test user.
//
//	@param name string
//	@param password string
//	@return int64  "user_id"
//	@return *models.User
//	@return string "token"
func registerTestUser(name string, password string) (int64, *models.User, string) {
	testUser := &models.User{
		Name:     name,
		Password: password,
	}

	rs := &services.RegisterService{}

	userID, _, _, token := rs.Register(testUser.Name, testUser.Password)

	return userID, testUser, token
}

// sendRequest sends a request to the server and
// returns the response recorder and response.
//
//	@param req *http.Request
//	@return *httptest.ResponseRecorder
//	@return interface{}
func sendRequest(req *http.Request) (*httptest.ResponseRecorder, interface{}) {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	res := selectResponseType(req)
	bytes, _ := io.ReadAll(w.Result().Body)
	json.Unmarshal(bytes, res)

	return w, res
}
