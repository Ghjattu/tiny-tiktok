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

// init() retrieves the environment variables, initializes the gin engine
// and registers a test user.
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
}

// selectResponseType selects the response type according to the request path.
//
//	@param req *http.Request
//	@return interface{}
func selectResponseType(req *http.Request) interface{} {
	pathSlice := strings.Split(req.URL.Path, "/")
	lastPath := pathSlice[len(pathSlice)-2]

	switch lastPath {
	case "feed":
		//  /douyin/feed/
		return &FeedResponse{}
	case "register":
		//  /douyin/user/register/
		return &RegisterResponse{}
	case "login":
		//  /douyin/user/login/
		return &LoginResponse{}
	case "user":
		//  /douyin/user/
		return &UserResponse{}
	case "action":
		//  /douyin/publish/action/
		return &Response{}
	case "list":
		//  /douyin/publish/list/
		return &PublishListResponse{}
	default:
		return nil
	}
}

// registerTestUser clears the database, and then registers a new test user
//
//	@param name string
//	@param password string
//	@return int64  "user_id"
//	@return *models.User
//	@return string "token"
func registerTestUser(name string, password string) (int64, *models.User, string) {
	models.InitDatabase(true)

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
