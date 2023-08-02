package router

import (
	"github.com/Ghjattu/tiny-tiktok/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	api := r.Group("/api")

	api.POST("/user/register", controllers.Register)
}
