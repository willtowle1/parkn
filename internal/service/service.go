package service

import (
	"context"
	"time"

	vision "cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/willtowle1/parkn/internal/common/errs"
	"github.com/willtowle1/parkn/internal/common/logger"
	"github.com/willtowle1/parkn/internal/model"
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

type IClient interface {
	FetchMedia(ctx context.Context, mediaUrl string) (*vision.Image, error)
}

type ParknService struct {
	logger        logger.Logger
	textExtractor ITextExtractor
	sniper        IDateSniper
	repository    IDal
	httpClient    IClient
}

func NewParknService(logger logger.Logger, textExtractor ITextExtractor, sniper IDateSniper, repository IDal, httpClient IClient) *ParknService {
	return &ParknService{
		logger:        logger,
		textExtractor: textExtractor,
		sniper:        sniper,
		repository:    repository,
		httpClient:    httpClient,
	}
}

// CreateParkn creates a parkn alert and returns the endDate
func (s *ParknService) CreateParkn(ctx context.Context, phoneNumber, mediaUrl string) (string, error) {

	image, err := s.httpClient.FetchMedia(ctx, mediaUrl)
	if err != nil {
		s.logger.Error(ctx, errCreatingParkn, err)
		return "", errs.WrapError(errCreatingParkn, err)
	}

	extractedText, err := s.textExtractor.ExtractTextFromImage(ctx, image)
	if err != nil {
		s.logger.Error(ctx, errCreatingParkn, err)
		return "", errs.WrapError(errCreatingParkn, err)
	}

	moveByDate, err := s.sniper.SnipeDate(ctx, extractedText)
	if err != nil {
		s.logger.Error(ctx, errCreatingParkn, err)
		return "", errs.WrapError(errCreatingParkn, err)
	}

	parknInput := model.Parkn{
		PhoneNumber: phoneNumber,
		MoveByDate:  moveByDate,
	}

	id, err := s.repository.CreateOne(ctx, parknInput)
	if err != nil {
		s.logger.Error(ctx, errCreatingParkn, err)
		return "", errs.WrapError(errCreatingParkn, err)
	}

	alertDate := fmtToString(moveByDate)
	s.logger.Info(ctx, msgCreateParknSuccess, "id", id, "alertDate", alertDate)

	return alertDate, nil
}

func fmtToString(t time.Time) string {
	return t.Format("01-02-2006")
}
