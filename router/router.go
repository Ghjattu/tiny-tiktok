package router

import (
	"github.com/Ghjattu/tiny-tiktok/controllers"
	"github.com/Ghjattu/tiny-tiktok/middleware/jwt"
	"github.com/Ghjattu/tiny-tiktok/middleware/parse"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	// Set static file path.
	r.Static("/static/videos", "../public/")

	api := r.Group("/douyin")
	api.Use(parse.ParseQueryParams())

	api.GET("/feed", jwt.AuthorizeFeed(), controllers.Feed)
	api.POST("/user/register/", controllers.Register)
	api.POST("/user/login/", controllers.Login)
	api.GET("/user/", jwt.AuthorizeGet(), controllers.GetUserByUserIDAndToken)
	api.POST("/publish/action/", jwt.AuthorizePost(), controllers.PublishNewVideo)
	api.GET("/publish/list/", jwt.AuthorizeGet(), controllers.GetPublishListByAuthorID)

	api.POST("/favorite/action/", jwt.AuthorizePost(), controllers.FavoriteAction)
	api.GET("/favorite/list/", jwt.AuthorizeGet(), controllers.GetFavoriteListByUserID)
	api.POST("/comment/action/", jwt.AuthorizePost(), controllers.CommentAction)
	api.GET("/comment/list/", jwt.AuthorizeGet(), controllers.CommentList)

	api.POST("/relation/action/", jwt.AuthorizePost(), controllers.FollowAction)
	api.GET("/relation/follow/list/", jwt.AuthorizeGet(), controllers.FollowingList)
	api.GET("/relation/follower/list/", jwt.AuthorizeGet(), controllers.FollowerList)
}
