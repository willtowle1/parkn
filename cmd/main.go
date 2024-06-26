package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/willtowle1/parkn/internal/common/logger"
	"github.com/willtowle1/parkn/internal/service"
)

func main() {

	ctx := context.Background()

	logger, err := logger.NewDefaultLogger("Debug")
	if err != nil {
		log.Fatal("failed to initialize logger ", err)
	}

	filename := "../data/IMG_5718.png"
	// filename := "../../data/IMG_5731.png"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	imageData, _ := io.ReadAll(f)

	b64Str := base64.StdEncoding.EncodeToString(imageData)

	imageClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	textExtractor := service.NewTextExtractor(logger, imageClient)
	dateSniper := service.NewDateSniper(logger)

	image, err := textExtractor.ConvertToVisionImage(ctx, b64Str)
	if err != nil {
		fmt.Println(err)
		return
	}

	text, err := textExtractor.ExtractTextFromImage(ctx, image)
	if err != nil {
		fmt.Println(err)
		return
	}

	endDate, err := dateSniper.SnipeDate(ctx, text)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(endDate.String())

	// client := openai.NewClient("sk-parkn-X4Se1KquYzKrWLFAYm5pT3BlbkFJcCSOfEqEVHUXMnDuIPIy")
	// gptClient := service.NewGPTClient(*logger, client)

	// date, err := gptClient.GetEndDate(context.Background(), text)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println(date)

	// extractedDate := extractFrequencyFromText(text)

	// fmt.Println(extractedDate)
}
