package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/willtowle1/parkn/internal/service"
)

func main() {
	ctx := context.Background()

	filename := "../data/IMG_5731.png"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	imageData, _ := io.ReadAll(f)

	b64Str := base64.StdEncoding.EncodeToString(imageData)
	fmt.Println(b64Str[:10])

	logger := log.New(io.Discard, "", log.Ldate)
	textExtractor := service.NewTextExtractor(*logger)

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

	fmt.Println(text)

}
