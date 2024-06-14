package service

import (
	"context"
	"log"

	"github.com/willtowle1/parkn/internal/model"
	"google.golang.org/genproto/googleapis/cloud/vision/v1"
)

type IParknRepository interface {
	CreateParkn(ctx context.Context, inputParkn model.Parkn) error
}

type ITextExtractor interface {
	ExtractTextFromImage(ctx context.Context, filename string) (string, error)
	ConvertToVisionImage(ctx context.Context, imageString string) (*vision.Image, error)
}

func NewParknService(logger log.Logger, textExtractor ITextExtractor, repository IParknRepository) *ParknService {
	return &ParknService{
		logger:        logger,
		textExtractor: textExtractor,
		repository:    repository,
	}
}

type ParknService struct {
	logger        log.Logger
	textExtractor ITextExtractor
	repository    IParknRepository
}

func (s *ParknService) CreateParkn(ctx context.Context) error {

	s.logger.Print("implement me!")

	return nil
}
