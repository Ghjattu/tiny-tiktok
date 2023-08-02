package main

import (
	"github.com/Ghjattu/tiny-tiktok/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	router.RegisterRouter(r)

	r.Run()
}
