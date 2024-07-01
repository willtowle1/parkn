package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/gin-gonic/gin"
	"github.com/willtowle1/parkn/internal/app"
	"github.com/willtowle1/parkn/internal/common/logger"
	"github.com/willtowle1/parkn/internal/config"
)

func main() {

	ctx := context.Background()

	logger, err := logger.NewDefaultLogger("Debug")
	if err != nil {
		log.Fatal("failed to initialize logger ", err)
	}

	config, err := config.Init(".env")
	if err != nil {
		logger.Error(ctx, "failed to get config", err)
		return
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	errs := make(chan error)

	imageClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get image client", err)
		return
	}

	mongoClient, err := app.InitDatabase(ctx, logger, errs, *config)
	if err != nil {
		logger.Error(ctx, "failed to get mongo client", err)
		return
	}
	database := mongoClient.Database(config.MongoDatabaseName)

	parknService := app.RegisterParknEndpoints(logger, router, imageClient, database)

	filename := "../data/IMG_5718.png"
	f, err := os.Open(filename)
	if err != nil {
		logger.Error(ctx, "failed to open file", err)
		return
	}
	defer f.Close()

	imageData, _ := io.ReadAll(f)

	b64Str := base64.StdEncoding.EncodeToString(imageData)

	endDate, err := parknService.CreateParkn(ctx, "+1 (314) 562-8484", b64Str)
	if err != nil {
		logger.Error(ctx, "failed to get endDate", err)
		return
	}
	fmt.Println(endDate)
}
