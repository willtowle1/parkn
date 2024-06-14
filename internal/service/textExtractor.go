package service

import (
	"context"
	"log"

	"google.golang.org/genproto/googleapis/cloud/vision/v1"
)

type TextExtractor struct {
	logger log.Logger
}

func NewTextExtractor(logger log.Logger) *TextExtractor {
	return &TextExtractor{
		logger: logger,
	}
}

// TODO: we won't want to use filename, will need to b64 encoding (?)
func (s *TextExtractor) ExtractTextFromImage(ctx context.Context, image *vision.Image) (string, error) {
	s.logger.Print("implement me!")
	return "", nil
}

func (s *TextExtractor) ConvertToVisionImage(b64Str string) (*vision.Image, error) {
	s.logger.Print("implement me!")
	return nil, nil
}
