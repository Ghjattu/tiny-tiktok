package rabbitmq

import (
	"log"

	"github.com/Ghjattu/tiny-tiktok/redis"
	"github.com/Ghjattu/tiny-tiktok/utils"
	"github.com/mitchellh/mapstructure"
)

func ProduceMessage(msgType, msgSubType, structName, key, field string, value interface{}) {
	message := &Message{
		Type:       msgType,
		SubType:    msgSubType,
		StructName: structName,
		Key:        key,
		Field:      field,
		Value:      value,
	}

	RedisMQ.Producer(message)
}

func ConsumeMessage(message *Message) {
	switch message.Type {
	case "Hash":
		switch message.SubType {
		case "Set":
			value := CacheStructSelector(message.StructName, message.Value)

			if err := redis.Rdb.HSet(redis.Ctx, message.Key, value).Err(); err != nil {
				log.Println("failed to hash set: ", err.Error())
			}
			redis.Rdb.Expire(redis.Ctx, message.Key, redis.RandomDay())
		case "Incr":
			incr, err := utils.ConvertInterfaceToInt64(message.Value)
			if err != nil {
				log.Println("Hash Incr error: ", err.Error())
			}

			if _, err := redis.HashIncrBy(message.Key, message.Field, incr); err != nil {
				log.Println("failed to hash incr by: ", err.Error())
			}
		}
	case "List":
		switch message.SubType {
		case "RPush":
			valueList, err := utils.ConvertInterfaceToInt64Slice(message.Value)
			if err != nil {
				log.Println("List RPush error: ", err.Error())
			}
			valueStrList, _ := utils.ConvertInt64ToString(valueList)

			if err := redis.Rdb.RPush(redis.Ctx, message.Key, valueStrList).Err(); err != nil {
				log.Println("failed to rpush: ", err.Error())
			}
			redis.Rdb.Expire(redis.Ctx, message.Key, redis.RandomDay())
		case "RPushX":
			valueList, err := utils.ConvertInterfaceToInt64Slice(message.Value)
			if err != nil {
				log.Println("List RPushX error: ", err.Error())
			}
			valueStrList, _ := utils.ConvertInt64ToString(valueList)

			if err := redis.Rdb.RPushX(redis.Ctx, message.Key, valueStrList).Err(); err != nil {
				log.Println("failed to rpushx: ", err.Error())
			}
		case "LRem":
			element, err := utils.ConvertInterfaceToInt64(message.Value)
			if err != nil {
				log.Println("List LRem error: ", err.Error())
			}

			if err := redis.Rdb.LRem(redis.Ctx, message.Key, 0, element).Err(); err != nil {
				log.Println("failed to lrem: ", err.Error())
			}
		}
	}
}

func CacheStructSelector(name string, messageValue interface{}) interface{} {
	switch name {
	case "VideoCache":
		videoCache := &redis.VideoCache{}
		if err := mapstructure.Decode(messageValue, videoCache); err != nil {
			log.Println("cache struct selector err: ", err.Error())
		}

		return videoCache
	case "CommentCache":
		commentCache := &redis.CommentCache{}
		if err := mapstructure.Decode(messageValue, commentCache); err != nil {
			log.Println("cache struct selector err: ", err.Error())
		}

		return commentCache
	case "UserCache":
		userCache := &redis.UserCache{}
		if err := mapstructure.Decode(messageValue, userCache); err != nil {
			log.Println("cache struct selector err: ", err.Error())
		}

		return userCache
	default:
		return nil
	}
}
