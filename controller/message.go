package controller

import (
	"douyin/models"
	"douyin/service"
	"douyin/utils/jwt"
	"douyin/utils/validator"
	"net/http"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

type MessageActionRequest struct {
	Token      string `query:"token" validate:"required"`
	ToUserID   string `query:"to_user_id" validate:"required"`
	ActionType string `query:"action_type" validate:"required,oneof=1 2"`
	Content    string `query:"content" validate:"required"`
}

type MessageChatRequest struct {
	Token      string `query:"token" validate:"required"`
	ToUserID   string `query:"to_user_id" validate:"required"`
	PreMsgTime string `query:"pre_msg_time" validate:"required"`
}

type ChatResponse struct {
	Response
	MessageList []models.MessageInfo `json:"message_list"`
	PreMsgTime  int64                `json:"pre_msg_time"`
}

func MessageAction(c *fiber.Ctx) error {
	request := MessageActionRequest{}
	emptyResponse := Response{}
	if err, httpErr := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return httpErr
	}
	var fromId uint
	if err, httpErr := jwt.JwtClient.AuthTokenValid(c, &emptyResponse, &fromId, request.Token); err != nil {
		return httpErr
	}
	toIdInt, _ := strconv.Atoi(request.ToUserID)
	toId := uint(toIdInt)
	content := request.Content

	service.AddMessage(toId, fromId, content)
	return c.Status(http.StatusOK).JSON(Response{StatusCode: 0, StatusMsg: "发送消息成功！"})
}

// MessageChat all users have same follow list
func MessageChat(c *fiber.Ctx) error {

	request := MessageChatRequest{}
	emptyResponse := Response{}
	if err, httpErr := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return httpErr
	}
	var fromId uint
	if err, httpErr := jwt.JwtClient.AuthTokenValid(c, &emptyResponse, &fromId, request.Token); err != nil {
		return httpErr
	}
	toIdnInt, _ := strconv.Atoi(request.ToUserID)
	toId := uint(toIdnInt)

	// 上次消息时间
	var preMsgTime int64
	preMsgTimeStr := request.PreMsgTime
	if preMsgTimeStr == "" {
		preMsgTime = 1546926630
	} else {
		preMsgTime, _ = strconv.ParseInt(preMsgTimeStr, 10, 64)
	}
	mids, err := service.GetMessagesIds(fromId, toId, &preMsgTime)
	nextPreMsgTime := preMsgTime
	// msgList, err := service.GetLatestMessageAfter(fromId, toId, preMsgTime)
	messages, err := service.GetMessagesByIds(mids)
	// 无消息
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(ChatResponse{
			Response:    Response{StatusCode: 1, StatusMsg: "no message"},
			MessageList: nil,
			PreMsgTime:  1546926630,
		})
	}
	messageInfos := make([]models.MessageInfo, len(mids))
	for i, m := range messages {
		messageInfos[i] = service.GenerateMessageInfo(&m)
	}
	// var nextPreMsgTime int64
	// if len(msgList) == 0 {
	// 	nextPreMsgTime = 1546926630
	// } else {
	// 	nextPreMsgTime = msgList[len(msgList)-1].CreateTime
	// }
	return c.Status(fiber.StatusOK).JSON(ChatResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "成功获取消息！"},
		MessageList: messageInfos,
		PreMsgTime:  nextPreMsgTime,
	})
}
