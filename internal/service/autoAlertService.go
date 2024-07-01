package service

import (
	"context"
	"strings"
	"time"

	"github.com/willtowle1/parkn/internal/common/logger"
)

const (
	errGettingParkns = "error while getting parkns to alert"
	errFailedToAlert = "error while alerting"
	errDeleteParkn   = "error while trying to delete parkn"

	msgAlertSuccessful  = "successfully sent alert"
	msgDeleteSuccessful = "successfully deleted parkn"
	msgAlertComplete    = "alert logic complete"
)

type VoiPService interface {
	SendAlert(ctx context.Context, phoneNumber string) error
}

type IAlertService interface {
	GetParknsToAlert(ctx context.Context, tomorrow time.Time) ([]string, error)
	DeleteParkn(ctx context.Context, phoneNumber string) (int64, error)
}

type AutoAlertService struct {
	logger  logger.Logger
	service IAlertService
	voip    VoiPService
}

func NewAutoAlertService(logger logger.Logger, service IAlertService, voip VoiPService) *AutoAlertService {
	return &AutoAlertService{
		logger:  logger,
		service: service,
		voip:    voip,
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
		err = s.voip.SendAlert(ctx, phoneNumber)
		if err != nil {
			// would be better to place into a separate queue to process later
			unsuccessful = append(unsuccessful, phoneNumber)
			s.logger.Error(ctx, errFailedToAlert, err, "phoneNumber", phoneNumber)
		} else {
			s.logger.Info(ctx, msgAlertSuccessful, "phoneNumber", phoneNumber)
			_, err = s.service.DeleteParkn(ctx, phoneNumber)
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
