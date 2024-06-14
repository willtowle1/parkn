package main

import (
	"context"
	"fmt"
	"log"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
)

func main() {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("failed to create client: %v ", err)
	}

	filename := "../data/IMG_5731.png"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		log.Fatal(err)
	}

	annotations, err := client.DetectTexts(ctx, image, nil, 10)
	if err != nil {
		log.Fatal(err)
	}

	if len(annotations) == 0 {
		log.Fatal("no annotations found")
	}

	fmt.Println("Text:")
	for _, annotation := range annotations {
		fmt.Println("Description: ", annotation.Description)
	}

}
