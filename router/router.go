package router

import (
	"github.com/Ghjattu/tiny-tiktok/controllers"
	"github.com/Ghjattu/tiny-tiktok/middleware/jwt"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	// Set static file path.
	r.Static("/static/videos", "../public/")

	api := r.Group("/douyin")

	api.GET("/feed", jwt.AuthorizationFeed(), controllers.Feed)
	api.POST("/user/register/", controllers.Register)
	api.POST("/user/login/", controllers.Login)
	api.GET("/user/", jwt.AuthorizationGet(), controllers.GetUserByUserIDAndToken)
	api.POST("/publish/action/", jwt.AuthorizationPost(), controllers.PublishNewVideo)
	api.GET("/publish/list/", jwt.AuthorizationGet(), controllers.GetPublishListByAuthorID)
}
