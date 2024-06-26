package app

import (
	"log"

	visionApi "cloud.google.com/go/vision/apiv1"
	"github.com/sashabaranov/go-openai"
	"github.com/willtowle1/parkn/internal/service"
)

func RegisterParknEndpoints(logger log.Logger, extractorClient *visionApi.ImageAnnotatorClient, openaiClient *openai.Client) {
	textExtractor := service.NewTextExtractor(logger, extractorClient)
	gptClient := service.NewGPTClient(logger, openaiClient)

	parknService := service.NewParknService(logger, textExtractor, gptClient)

}
