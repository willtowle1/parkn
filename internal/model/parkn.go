package model

import "time"

type Parkn struct {
	PhoneNumber string    `json:"phoneNumber"`
	MoveByDate  time.Time `json:"moveByDate"`
}
