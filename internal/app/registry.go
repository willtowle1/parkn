package app

import (
	visionApi "cloud.google.com/go/vision/apiv1"
	"github.com/gin-gonic/gin"
	"github.com/willtowle1/parkn/internal/common/logger"
	"github.com/willtowle1/parkn/internal/controller"
	"github.com/willtowle1/parkn/internal/dal"
	"github.com/willtowle1/parkn/internal/model"
	"github.com/willtowle1/parkn/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterParknEndpoints(logger logger.Logger, router gin.IRouter, extractorClient *visionApi.ImageAnnotatorClient, database *mongo.Database) *service.ParknService {

	parknCollection := database.Collection("parkns")
	parknRepository := dal.NewRepository[model.Parkn](logger, *parknCollection)

	parknTextExtractor := service.NewTextExtractor(logger, extractorClient)
	parknDateSniper := service.NewDateSniper(logger)

	parknService := service.NewParknService(logger, parknTextExtractor, parknDateSniper, parknRepository)

	parknController := controller.NewController(logger, parknService)

	apiRouter := router.Group("/api")
	parknController.RegisterRoutes(apiRouter)

	return parknService
}
