package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/Ghjattu/tiny-tiktok/redis"
	"github.com/Ghjattu/tiny-tiktok/services"
	"github.com/Ghjattu/tiny-tiktok/utils"
	"github.com/gin-gonic/gin"
)

type MessageChatResponse struct {
	Response
	MessageList []models.MessageDetail `json:"message_list"`
}

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

func MessageChat(c *gin.Context) {
	receiverID := c.GetInt64("to_user_id")
	currentUserID := c.GetInt64("current_user_id")
	preMsgTimeInt := c.GetInt64("pre_msg_time")
	preMsgTime := time.Unix(preMsgTimeInt, 0)

	if preMsgTime.Year() > 9999 {
		// Get last message time from redis.
		lastMsgTimeKey := redis.LastMsgTimeKey + strconv.FormatInt(currentUserID, 10) + ":" +
			strconv.FormatInt(receiverID, 10)
		preMsgTimeStr := redis.Rdb.Get(redis.Ctx, lastMsgTimeKey).Val()

		_, _, preMsgTimeInt = utils.ParseInt64(preMsgTimeStr)
		preMsgTime = time.Unix(preMsgTimeInt, 0)
	}

	ms := &services.MessageService{}
	statusCode, statusMsg, messages := ms.GetMessageList(currentUserID, receiverID, preMsgTime)

	c.JSON(http.StatusOK, MessageChatResponse{
		Response: Response{
			StatusCode: statusCode,
			StatusMsg:  statusMsg,
		},
		MessageList: messages,
	})
}
