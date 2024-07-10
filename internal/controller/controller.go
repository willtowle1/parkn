package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go/twiml"
	"github.com/willtowle1/parkn/internal/common/logger"
)

const (
	msgCreateParknSuccess = "parkn alert created successfully"

	errCreateParkn          = "error while creating parkn alert"
	errMissingPhoneNumber   = "no phone number found in context"
	errMissingMedia         = "no media found in message"
	errMissingImageEncoding = "no image encoding found in context"
)

type IService interface {
	CreateParkn(ctx context.Context, phoneNumber, mediaUrl string) (string, error)
}

type Controller struct {
	logger  logger.Logger
	service IService
}

func NewController(logger logger.Logger, service IService) *Controller {
	return &Controller{
		logger:  logger,
		service: service,
	}
}

func (c *Controller) RegisterRoutes(router gin.IRouter) {
	route := router.Group("/v1")
	route.Handle(http.MethodPost, "/parkn/sms", c.createParkn)
}

func (c *Controller) createParkn(ctx *gin.Context) {

	ctx.Header("Content-Type", "text/xml")

	phoneNumber := ctx.PostForm("From")
	if len(phoneNumber) == 0 {
		err := errors.New(errMissingPhoneNumber)
		c.logger.Error(ctx, errCreateParkn, err)
		message := c.createErrorMessage(errCreateParkn, errMissingPhoneNumber)
		ctx.String(http.StatusBadRequest, message)
		return
	}

	mediaUrl := ctx.PostForm("MediaUrl0")
	if len(mediaUrl) == 0 {
		err := errors.New(errMissingMedia)
		c.logger.Error(ctx, errCreateParkn, err)
		message := c.createErrorMessage(errCreateParkn, errMissingMedia)
		ctx.String(http.StatusBadRequest, message)
		return
	}

	moveByDate, err := c.service.CreateParkn(ctx, phoneNumber, mediaUrl)

	if err != nil {
		c.logger.Error(ctx, errCreateParkn, err)
		message := c.createErrorMessage(errCreateParkn, err.Error())
		ctx.String(http.StatusInternalServerError, message)
		return
	}

	message := c.createSuccessMessage(msgCreateParknSuccess, moveByDate)

	c.logger.Info(ctx, msgCreateParknSuccess, "moveByDate", moveByDate)
	ctx.String(http.StatusOK, message)
}

func (c *Controller) createErrorMessage(msg, errString string) string {
	message := &twiml.MessagingMessage{
		Body: fmt.Sprintf("Error - %s: %s", msg, errString),
	}
	res, _ := twiml.Messages([]twiml.Element{message})
	return res
}

func (c *Controller) createSuccessMessage(msg, moveByDate string) string {
	message := &twiml.MessagingMessage{
		Body: fmt.Sprintf("Success - %s. You will be alerted to move your car by %s", msg, moveByDate),
	}
	res, _ := twiml.Messages([]twiml.Element{message})
	return res
}
