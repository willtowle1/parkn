package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sashabaranov/go-openai"
)

var (
	prompt1 = "Given this text, provide me with the next instance of the event in the format \"MM-DD-YYYY:HH:MM:SS\". Today is %s %d, %d. Please provide only the date in the given format and no other text. \"%s\""
)

type GPTClient struct {
	logger log.Logger
	client openai.Client
}

func NewGPTClient(logger log.Logger, client openai.Client) *GPTClient {
	return &GPTClient{
		logger: logger,
		client: client,
	}
}

// GetEndDate uses pre-defined prompts to extract date from text input using go-openai library
func (c *GPTClient) GetEndDate(ctx context.Context, inputText string) (time.Time, error) {
	resp, err := c.client.
}

func (c *GPTClient) formatPrompt(inputText string) string {
	year, month, day := time.Now().Date()
	return fmt.Sprintf(prompt1, month, day, year, inputText)
}

func (c *GPTClient) formatResponse(outputText string) time.Time {
	f.Fatalf("implement me!")
	return time.Now()
}
