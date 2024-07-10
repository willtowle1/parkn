package service

import (
	"context"
	"time"

	"github.com/willtowle1/parkn/internal/common/errs"
	"github.com/willtowle1/parkn/internal/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	errGetParknsToAlert = "error while getting parkns to alert"
	errDeletingParkn    = "error while deleting parkn"
	errNoParknDeleted   = "delete count of zero while deleting parkn"
)

type AlertService struct {
	logger     logger.Logger
	repository IDal
}

func NewAlertService(logger logger.Logger, repository IDal) *AlertService {
	return &AlertService{
		logger:     logger,
		repository: repository,
	}
}

func (s *AlertService) GetParknsToAlert(ctx context.Context, tomorrow time.Time) ([]string, error) {

	filter := bson.D{
		{Key: "moveByDate", Value: bson.D{
			{Key: "$lte", Value: primitive.NewDateTimeFromTime(tomorrow)},
		}},
	}

	phoneNumbers := make([]string, 0)
	parkns, err := s.repository.Get(ctx, filter)
	if err != nil {
		return phoneNumbers, errs.WrapError(errGetParknsToAlert, err)
	}

	for _, parkn := range parkns {
		phoneNumbers = append(phoneNumbers, parkn.PhoneNumber)
	}

	return phoneNumbers, nil
}

func (s *AlertService) DeleteParkn(ctx context.Context, phoneNumber string) error {
	filter := bson.D{
		{Key: "phoneNumber", Value: phoneNumber},
	}

	deleteCount, err := s.repository.DeleteOne(ctx, filter)
	if err != nil {
		return errs.WrapError(errDeleteParkn, err)
	}
	if deleteCount == 0 {
		return errs.WrapError(errNoParknDeleted, err)
	}

	return nil
}
