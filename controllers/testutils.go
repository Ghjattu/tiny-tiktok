// Description:
// This testutils.go file contains some functions and variables
// that are used in test files in the controllers package.

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/Ghjattu/tiny-tiktok/middleware/jwt"
	"github.com/Ghjattu/tiny-tiktok/middleware/parse"
	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/gin-gonic/gin"
)

var (
	r *gin.Engine

	userID    int64
	userIDStr string
	token     string
)

func setup() {
	models.InitDatabase(true)

	// Register a test user.
	userID, _, token = registerTestUser("test", "123456")
	userIDStr = fmt.Sprintf("%d", userID)

	r = gin.Default()
	r.Use(parse.ParseQueryParams())
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
	r.GET("/douyin/relation/follow/list/", jwt.AuthorizeGet(), FollowingList)
	r.GET("/douyin/relation/follower/list/", jwt.AuthorizeGet(), FollowerList)
	r.GET("/douyin/relation/friend/list/", jwt.AuthorizeGet(), FriendList)
	r.POST("/douyin/message/action/", jwt.AuthorizePost(), MessageAction)
	r.GET("/douyin/message/chat/", jwt.AuthorizeGet(), MessageChat)
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
	case "/relation/follow/list/":
		return &UserListResponse{}
	case "/relation/follower/list/":
		return &UserListResponse{}
	case "/relation/friend/list/":
		return &FriendListResponse{}
	case "/message/chat/":
		return &MessageChatResponse{}
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

// constructTestForm constructs a test form with a test file and form fields.
//
//	@param formFields map[string]string
//	@return *bytes.Buffer
//	@return *multipart.Writer
//	@return error
func constructTestForm(formFields map[string]string) (*bytes.Buffer, *multipart.Writer, error) {
	// Read the test video.
	file, err := os.Open("../data/bear.mp4")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// Construct form data.
	form := bytes.NewBuffer([]byte(""))
	writer := multipart.NewWriter(form)
	defer writer.Close()

	// Add form fields.
	for key, value := range formFields {
		writer.WriteField(key, value)
	}

	// Add form file.
	part, err := writer.CreateFormFile("data", "bear.mp4")
	if err != nil {
		return nil, nil, err
	}

	// Copy file data to form file.
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, nil, err
	}

	return form, writer, nil
}
