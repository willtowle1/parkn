package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twilio/twilio-go/twiml"
)

func main() {

	router := gin.Default()

	router.POST("/sms", func(context *gin.Context) {

		incomingMsg := context.PostForm("Body")
		incomingPhoneNumber := context.PostForm("From")
		mediaUrl := context.PostForm("MediaUrl0")
		fmt.Println(incomingMsg, incomingPhoneNumber, mediaUrl)

		// https://help.twilio.com/articles/223181368
		// https://www.twilio.com/docs/messaging/api/media-resource#fetch-a-media
		// look into setting env vars properly (TWILIO_ACCOUNT_SID & TWILIO_AUTH_TOKEN)

		message := &twiml.MessagingMessage{
			Body: "Hello from within!",
		}

		res, err := twiml.Messages([]twiml.Element{message})
		if err != nil {
			fmt.Println(err)
			return
		}
		context.Header("Content-Type", "text/xml")
		context.String(http.StatusOK, res)

	})

	router.Run(":3000")

	// ctx := context.Background()

	// logger, err := logger.NewDefaultLogger("Debug")
	// if err != nil {
	// 	log.Fatal("failed to initialize logger ", err)
	// }

	// config, err := config.Init(".env")
	// if err != nil {
	// 	logger.Error(ctx, "failed to get config", err)
	// 	return
	// }

	// gin.SetMode(gin.ReleaseMode)
	// router := gin.New()

	// errs := make(chan error)

	// imageClient, err := vision.NewImageAnnotatorClient(ctx)
	// if err != nil {
	// 	logger.Error(ctx, "failed to get image client", err)
	// 	return
	// }

	// mongoClient, err := app.InitDatabase(ctx, logger, errs, *config)
	// if err != nil {
	// 	logger.Error(ctx, "failed to get mongo client", err)
	// 	return
	// }
	// database := mongoClient.Database(config.MongoDatabaseName)

	// parknService := app.RegisterParknEndpoints(logger, router, imageClient, database)

	// filename := "../data/IMG_5718.png"
	// f, err := os.Open(filename)
	// if err != nil {
	// 	logger.Error(ctx, "failed to open file", err)
	// 	return
	// }
	// defer f.Close()

	// imageData, _ := io.ReadAll(f)

	// b64Str := base64.StdEncoding.EncodeToString(imageData)

	// endDate, err := parknService.CreateParkn(ctx, "+1 (314) 562-8484", b64Str)
	// if err != nil {
	// 	logger.Error(ctx, "failed to get endDate", err)
	// 	return
	// }
	// fmt.Println(endDate)
}
