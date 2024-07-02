package client

import (
	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go"
	"github.com/willtowle1/parkn/internal/common/logger"
)

type Client struct {
	logger logger.Logger
	client twilio.RestClient
}

func NewTwilioClient(logger logger.Logger, client twilio.RestClient) *Client {
	return &Client{
		logger: logger,
		client: client,
	}
}

// incoming messages should be handled in controller
// this package should handle sending outbound messages

func (c *Client) HandleIncomingMessages(router gin.IRouter) {
	panic("implement me!")
}
