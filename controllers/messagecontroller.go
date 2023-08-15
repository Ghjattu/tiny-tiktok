package controllers

import (
	"net/http"

	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/gin-gonic/gin"
)

func MessageAction(c *gin.Context) {
	receiverID := c.GetInt64("to_user_id")
	actionType := c.GetInt64("action_type")
	content := c.Query("content")
	currentUserID := c.GetInt64("current_user_id")

	statusCode := int32(1)
	statusMsg := "action type is invalid"

	ms := &services.MessageService{}
	if actionType == 1 {
		statusCode, statusMsg = ms.CreateNewMessage(currentUserID, receiverID, content)
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	})
}
