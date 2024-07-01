package service

import (
	"context"
	"regexp"
	"time"

	"github.com/willtowle1/parkn/internal/common/errs"
	"github.com/willtowle1/parkn/internal/common/logger"
	"github.com/willtowle1/parkn/internal/model"
	"google.golang.org/genproto/googleapis/cloud/vision/v1"
)

const (
	errCreatingParkn = "failed to create parkn"

	msgCreateParknSuccess = "successfully created parkn alert"
)

type IDal interface {
	CreateOne(ctx context.Context, input model.Parkn) (string, error)
	Get(ctx context.Context, filter interface{}) ([]model.Parkn, error)
	DeleteOne(ctx context.Context, filter interface{}) (int64, error)
}

type ITextExtractor interface {
	ExtractTextFromImage(ctx context.Context, image *vision.Image) (string, error)
	ConvertToVisionImage(ctx context.Context, imageString string) (*vision.Image, error)
}

type IDateSniper interface {
	SnipeDate(ctx context.Context, str string) (time.Time, error)
}

func NewParknService(logger logger.Logger, textExtractor ITextExtractor, sniper IDateSniper, repository IDal) *ParknService {
	return &ParknService{
		logger:        logger,
		textExtractor: textExtractor,
		sniper:        sniper,
		repository:    repository,
	}
}

type ParknService struct {
	logger        logger.Logger
	textExtractor ITextExtractor
	sniper        IDateSniper
	repository    IDal
}

// CreateParkn creates a parkn alert and returns the endDate
func (s *ParknService) CreateParkn(ctx context.Context, phoneNumber, imageEncoding string) (string, error) {

	image, err := s.textExtractor.ConvertToVisionImage(ctx, imageEncoding)
	if err != nil {
		s.logger.Error(ctx, errCreatingParkn, err)
		return "", errs.WrapError(errCreatingParkn, err)
	}

	extractedText, err := s.textExtractor.ExtractTextFromImage(ctx, image)
	if err != nil {
		s.logger.Error(ctx, errCreatingParkn, err)
		return "", errs.WrapError(errCreatingParkn, err)
	}

	endDate, _ := s.sniper.SnipeDate(ctx, extractedText)

	parknInput := model.Parkn{
		PhoneNumber: s.cleanPhoneNumber(phoneNumber),
		MoveByDate:  endDate,
	}

	id, err := s.repository.CreateOne(ctx, parknInput)
	if err != nil {
		s.logger.Error(ctx, errCreatingParkn, err)
		return "", errs.WrapError(errCreatingParkn, err)
	}

	alertDate := endDate.String()
	s.logger.Info(ctx, msgCreateParknSuccess, "id", id, "alertDate", alertDate)

	return alertDate, nil
}

func truncateToMinute(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
}

func (s *ParknService) cleanPhoneNumber(str string) string {
	r := regexp.MustCompile(`\D`)
	return r.ReplaceAllString(str, "")
}
