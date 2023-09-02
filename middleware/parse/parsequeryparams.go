package parse

import (
	"net/http"

	"github.com/Ghjattu/tiny-tiktok/utils"
	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func ParseQueryParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		validParams := map[string]string{
			"user_id":      "int",
			"video_id":     "int",
			"action_type":  "int",
			"comment_text": "string",
			"comment_id":   "int",
			"to_user_id":   "int",
			"content":      "string",
			"pre_msg_time": "int",
		}

		queryParams := c.Request.URL.Query()
		queryValues := make(map[string]interface{})

		for key, valueList := range queryParams {
			// If the key is not valid, ignore it.
			if _, ok := validParams[key]; !ok {
				continue
			}
			// Get the first value of the list.
			valueStr := valueList[0]
			// If the type of the key is string, just save it.
			if t := validParams[key]; t == "string" {
				queryValues[key] = valueStr
				continue
			}
			// If the type of the key is int, parse it.
			statusCode, statusMsg, valueInt := utils.ParseInt64(valueStr)
			if statusCode != 0 {
				c.JSON(http.StatusOK, Response{
					StatusCode: statusCode,
					StatusMsg:  statusMsg,
				})
				c.Abort()
				return
			}
			queryValues[key] = valueInt
		}

		// If all the params are valid, set them to the context.
		for key, value := range queryValues {
			if t := validParams[key]; t == "string" {
				c.Set(key, value.(string))
			} else {
				c.Set(key, value.(int64))
			}
		}

		c.Next()
	}
}
