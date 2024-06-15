package service

import (
	"context"
	"log"
	"time"

	"github.com/willtowle1/parkn/internal/common/errs"
	"github.com/willtowle1/parkn/internal/model"
	"google.golang.org/genproto/googleapis/cloud/vision/v1"
)

const (
	errCreatingParkn = "failed to create parkn"
)

type IParknRepository interface {
	CreateParkn(ctx context.Context, inputParkn model.Parkn) error
}

type ITextExtractor interface {
	ExtractTextFromImage(ctx context.Context, image *vision.Image) (string, error)
	ConvertToVisionImage(ctx context.Context, imageString string) (*vision.Image, error)
}

type IGPTClient interface {
	GetEndDate(ctx context.Context, inputText string) (time.Time, error)
}

func NewParknService(logger log.Logger, textExtractor ITextExtractor, gptClient IGPTClient, repository IParknRepository) *ParknService {
	return &ParknService{
		logger:        logger,
		textExtractor: textExtractor,
		gptClient:     gptClient,
		repository:    repository,
	}
}

type ParknService struct {
	logger        log.Logger
	textExtractor ITextExtractor
	gptClient     IGPTClient
	repository    IParknRepository
}

// CreateParkn creates a parkn alert and returns the endDate
func (s *ParknService) CreateParkn(ctx context.Context, phoneNumber, imageEncoding string) (string, error) {

	image, err := s.textExtractor.ConvertToVisionImage(ctx, imageEncoding)
	if err != nil {
		return "", errs.WrapError(errCreatingParkn, err)
	}

	extractedText, err := s.textExtractor.ExtractTextFromImage(ctx, image)
	if err != nil {
		return "", errs.WrapError(errCreatingParkn, err)
	}

	endDate, err := s.gptClient.GetEndDate(ctx, extractedText)
	if err != nil {
		return "", errs.WrapError(errCreatingParkn, err)
	}

	parkn := model.Parkn{
		PhoneNumber: phoneNumber,
		StartDate:   time.Now().UTC().String(),
		EndDate:     endDate.String(),
	}

	err = s.repository.CreateParkn(ctx, parkn)
	if err != nil {
		return "", errs.WrapError(errCreatingParkn, err)
	}

	return endDate.String(), nil
}
