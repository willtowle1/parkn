package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/willtowle1/parkn/internal/common/logger"
)

const (
	msgCreateParknSuccess = "parkn alert created successfully"

	errCreateParkn = "error while creating parkn alert"

	phoneNumberKey   = "phoneNumber"
	imageEncodingKey = "imageEncoding"
)

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

}
