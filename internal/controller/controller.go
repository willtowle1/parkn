package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/willtowle1/parkn/internal/common/errs"
	"github.com/willtowle1/parkn/internal/common/logger"
)

const (
	msgCreateParknSuccess = "parkn alert created successfully"

	errCreateParkn          = "error while creating parkn alert"
	errMissingPhoneNumber   = "no phone number found in context"
	errMissingImageEncoding = "no image encoding found in context"

	phoneNumberKey   = "phoneNumber"
	imageEncodingKey = "imageEncoding"
)

type successfulResponse struct {
	AlertDate string `json:"alertDate"`
}

type IService interface {
	CreateParkn(ctx context.Context, phoneNumber, imageEncoding string) (string, error)
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
	route.Handle(http.MethodPost, "/parkn", c.createParkn)
}

func (c *Controller) createParkn(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	phoneNumber, ok := reqCtx.Value(phoneNumberKey).(string)
	if !ok {
		err := errors.New(errMissingPhoneNumber)
		c.logger.Error(ctx, errCreateParkn, err)
		apiErr := errs.NewApiError(http.StatusBadRequest, "BAD_REQUEST", errCreateParkn, "error", err.Error())
		ctx.AbortWithStatusJSON(apiErr.Status, apiErr)
		return
	}

	imageEncoding, ok := reqCtx.Value(imageEncodingKey).(string)
	if !ok {
		err := errors.New(errMissingImageEncoding)
		c.logger.Error(ctx, errCreateParkn, err)
		apiErr := errs.NewApiError(http.StatusBadRequest, "BAD_REQUEST", errCreateParkn, "error", err.Error())
		ctx.AbortWithStatusJSON(apiErr.Status, apiErr)
		return
	}

	alertDate, err := c.service.CreateParkn(ctx, phoneNumber, imageEncoding)

	if err != nil {
		c.logger.Error(ctx, errCreateParkn, err)
		apiErr := errs.NewApiError(http.StatusInternalServerError, "INTERNAL", errCreateParkn, "error", err.Error())
		ctx.AbortWithStatusJSON(apiErr.Status, apiErr)
		return
	}

	resp := successfulResponse{
		AlertDate: alertDate,
	}
	c.logger.Info(ctx, msgCreateParknSuccess, "alertDate", alertDate)
	ctx.JSON(http.StatusOK, resp)
}
