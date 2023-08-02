package main

import (
	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	router.RegisterRouter(r)

	Initialize()

	r.Run()
}

func Initialize() {
	models.InitDatabase()
}
