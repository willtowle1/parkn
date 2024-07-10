package service

import (
	"context"
	"strings"
	"time"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/willtowle1/parkn/internal/common/logger"
)

const (
	errGettingParkns = "error while getting parkns to alert"
	errFailedToAlert = "error while alerting"
	errDeleteParkn   = "error while trying to delete parkn"

	msgAlertSuccessful  = "successfully sent alert"
	msgDeleteSuccessful = "successfully deleted parkn"
	msgAlertComplete    = "alert logic complete"

	alertMsg = "Move your car by tomorrow!"
)

type IAlertService interface {
	GetParknsToAlert(ctx context.Context, tomorrow time.Time) ([]string, error)
	DeleteParkn(ctx context.Context, phoneNumber string) error
}

type AutoAlertService struct {
	logger       logger.Logger
	service      IAlertService
	twilio       *twilio.RestClient
	twilioNumber string
}

func NewAutoAlertService(logger logger.Logger, service IAlertService, twilio *twilio.RestClient, twilioNumber string) *AutoAlertService {
	return &AutoAlertService{
		logger:       logger,
		service:      service,
		twilio:       twilio,
		twilioNumber: twilioNumber,
	}
}

func (s *AutoAlertService) Alert(ctx context.Context) {
	loc, _ := time.LoadLocation("EST")
	tomorrow := time.Now().In(loc).Add(24 * time.Hour)

	toAlert, err := s.service.GetParknsToAlert(ctx, tomorrow)
	if err != nil {
		s.logger.Error(ctx, errGettingParkns, err)
		return
	}

	successful := make([]string, 0)
	unsuccessful := make([]string, 0)
	for _, phoneNumber := range toAlert {
		err = s.sendAlert(phoneNumber)
		if err != nil {
			// TODO: would be better to place into a separate queue to process later
			unsuccessful = append(unsuccessful, phoneNumber)
			s.logger.Error(ctx, errFailedToAlert, err, "phoneNumber", phoneNumber)
		} else {
			s.logger.Info(ctx, msgAlertSuccessful, "phoneNumber", phoneNumber)
			err = s.service.DeleteParkn(ctx, phoneNumber)
			if err != nil {
				unsuccessful = append(unsuccessful, phoneNumber)
				s.logger.Error(ctx, errDeleteParkn, err, "phoneNumber", phoneNumber)
			} else {
				s.logger.Info(ctx, msgDeleteSuccessful, "phoneNumber", phoneNumber)
				successful = append(successful, phoneNumber)
			}
		}
	}

	s.logger.Info(ctx, msgAlertComplete, "successful", strings.Join(successful, ", "), "unsuccessful", strings.Join(unsuccessful, ", "))
}

func (s *AutoAlertService) sendAlert(phoneNumber string) error {

	params := &twilioApi.CreateMessageParams{
		To:   s.strPtr(phoneNumber),
		From: s.strPtr(s.twilioNumber),
		Body: s.strPtr(alertMsg),
	}

	_, err := s.twilio.Api.CreateMessage(params)

	return err
}

func (s *AutoAlertService) strPtr(str string) *string {
	return &str
}
