package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/twilio/twilio-go"
	"github.com/willtowle1/parkn/internal/app"
	"github.com/willtowle1/parkn/internal/common/logger"
	"github.com/willtowle1/parkn/internal/config"
)

func main() {

	ctx := context.Background()

	config, err := config.Init(".env")
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	logger, err := logger.NewDefaultLogger(config.LogLevel)
	if err != nil {
		log.Fatalf("failed to get new logger: %s", err)
	}

	errs := make(chan error)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	extractorClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get image client", err)
		os.Exit(1)
	}

	mongoClient, err := app.InitDatabase(ctx, logger, errs, *config)
	if err != nil {
		logger.Error(ctx, "failed to get mongo client", err)
		os.Exit(1)
	}
	database := mongoClient.Database(config.MongoDatabaseName)

	twilioCreds := twilio.ClientParams{
		Username: config.TwilioSID,
		Password: config.TwilioToken,
	}
	twilioClient := twilio.NewRestClientWithParams(twilioCreds)
	app.RegisterParknEndpoints(logger, router, extractorClient, database, twilioCreds)
	autoAlertService := app.RegisterAutoAlertService(logger, database, twilioClient, config.TwilioNumber)

	scheduler := gocron.NewScheduler(time.UTC)
	_, err = scheduler.Every(config.AutoAlertPeriod).Minute().Do(autoAlertService.Alert, ctx)
	if err != nil {
		logger.Error(ctx, "failed to start auto service", err)
	}

	mainApp := app.NewApp(logger, &http.Server{
		Addr:    config.ServerAddress,
		Handler: router,
	})

	mainApp.Start(ctx, errs, config.ServerAddress)

	scheduler.StartAsync()

	app.WaitForTermination(ctx, logger, errs)

	scheduler.Stop()

	err = mainApp.Shutdown(ctx, time.Duration(config.TerminationGracePeriod)*time.Second)
	if err != nil {
		logger.Error(ctx, "error while shutting down server", err)
	} else {
		logger.Info(ctx, "server terminated successfully")
	}
}
