package logger

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	levelMap = map[string]int{
		"Error": 1,
		"Info":  2,
		"Debug": 3,
	}
)

type Logger interface {
	Error(ctx context.Context, message string, err error, keyValues ...interface{})
	Info(ctx context.Context, message string, keyValues ...interface{})
	Debug(ctx context.Context, message string, keyValues ...interface{})
}

type DefaultLogger struct {
	level int
}

func NewDefaultLogger(level string) (*DefaultLogger, error) {
	levelInt, exists := levelMap[level]
	if !exists {
		return nil, errors.New("invalid level")
	}
	return &DefaultLogger{
		level: levelInt,
	}, nil
}

func (l *DefaultLogger) Error(ctx context.Context, message string, err error, keyValues ...interface{}) {
	fmt.Println(l.fmtLog("ERROR", message, err, keyValues...))
}

func (l *DefaultLogger) Info(ctx context.Context, message string, keyValues ...interface{}) {
	if l.level < levelMap["Info"] {
		return
	}
	fmt.Println(l.fmtLog("INFO", message, nil, keyValues...))
}

func (l *DefaultLogger) Debug(ctx context.Context, message string, keyValues ...interface{}) {
	if l.level < levelMap["Debug"] {
		return
	}
	fmt.Println(l.fmtLog("DEBUG", message, nil, keyValues...))
}

func (l *DefaultLogger) fmtLog(logLevel, message string, err error, keyValues ...interface{}) string {
	str := ""
	i := 0
	for i < len(keyValues) {
		key := keyValues[i]
		var value interface{}
		if i+1 < len(keyValues) {
			value = keyValues[i+1]
		}
		str += fmt.Sprintf("%s: %s... ", key, value)
		i += 2
	}

	loc, _ := time.LoadLocation("EST")

	timeNow := truncateToMinuteString(time.Now().In(loc))

	if err != nil {
		return fmt.Sprintf("%s... %s -> %s: %s... %s", timeNow, logLevel, message, err.Error(), str)
	}

	return fmt.Sprintf("%s... %s -> %s... %s", timeNow, logLevel, message, str)
}

func truncateToMinuteString(t time.Time) string {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location()).Format("01-02-2006 15:04:00")
}
