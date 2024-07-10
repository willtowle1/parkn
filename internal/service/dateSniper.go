package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/teambition/rrule-go"
	"github.com/willtowle1/parkn/internal/common/errs"
	"github.com/willtowle1/parkn/internal/common/logger"
)

const (
	errSnipingDate           = "error while finding date"
	errNoFrequencyFound      = "no frequency detected in extracted text from image"
	errCalculatingOccurrence = "error while calculating next occurrence"
)

type frequency struct {
	dayOfWeek        int
	firstOccOfMonth  int
	secondOccOfMonth int
}

var (
	possibleFrequencies = map[string]frequency{
		"1ST & 3RD MONDAY":    {1, 1, 3},
		"2ND & 4TH MONDAY":    {1, 2, 4},
		"1ST & 3RD TUESDAY":   {2, 1, 3},
		"2ND & 4TH TUESDAY":   {2, 2, 4},
		"1ST & 3RD WEDNESDAY": {3, 1, 3},
		"2ND & 4TH WEDNESDAY": {3, 2, 4},
		"1ST & 3RD THURSDAY":  {4, 1, 3},
		"2ND & 4TH THURSDAY":  {4, 2, 4},
		"1ST & 3RD FRIDAY":    {5, 1, 3},
		"2ND & 4TH FRIDAY":    {5, 2, 4},
		"1ST & 3RD SATURDAY":  {6, 1, 3},
		"2ND & 4TH SATURDAY":  {6, 2, 4},
		"1ST & 3RD SUNDAY":    {7, 1, 3},
		"2ND & 4TH SUNDAY":    {7, 2, 4},
	}

	weekdayToRule = map[int]interface{}{
		1: rrule.MO,
		2: rrule.TU,
		3: rrule.WE,
		4: rrule.TH,
		5: rrule.FR,
		6: rrule.SA,
		7: rrule.SU,
	}
)

type DateSniper struct {
	logger logger.Logger
}

func NewDateSniper(logger logger.Logger) *DateSniper {
	return &DateSniper{
		logger: logger,
	}
}

// SnipeDate takes extracted image text and finds next street sweeping occurrence
func (d *DateSniper) SnipeDate(ctx context.Context, str string) (time.Time, error) {

	strArr := strings.Split(strings.ToUpper(str), "\n")

	freq, found := d.getFreq(strArr)
	if !found {
		err := errors.New(errNoFrequencyFound)
		d.logger.Error(ctx, errSnipingDate, err)
		return time.Time{}, errs.WrapError(errSnipingDate, err)
	}

	nextOccurrence, err := d.findNextOccurrence(freq)
	if err != nil {
		d.logger.Error(ctx, errSnipingDate, err)
		return time.Time{}, errs.WrapError(errSnipingDate, err)
	}
	return nextOccurrence, nil
}

func (d *DateSniper) findNextOccurrence(freq frequency) (time.Time, error) {

	loc, err := time.LoadLocation("EST")
	if err != nil {
		return time.Time{}, err
	}

	rruleWeekday := weekdayToRule[freq.dayOfWeek]
	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, loc)

	rule, err := rrule.NewRRule(rrule.ROption{
		Freq:      rrule.MONTHLY,
		Dtstart:   startDate,
		Byweekday: []rrule.Weekday{rruleWeekday.(rrule.Weekday)},
		Bysetpos:  []int{freq.firstOccOfMonth, freq.secondOccOfMonth},
	})

	if err != nil {
		return time.Time{}, err
	}

	today := time.Now().Add(time.Hour * -24).In(loc)
	nextOccurrence := rule.After(today, true).In(loc)
	return d.truncateToDay(nextOccurrence), nil
}

func (d *DateSniper) getFreq(strArr []string) (frequency, bool) {
	for _, str := range strArr {
		if freq, exists := possibleFrequencies[str]; exists {
			return freq, true
		}
	}
	return frequency{}, false
}

func (d *DateSniper) truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
