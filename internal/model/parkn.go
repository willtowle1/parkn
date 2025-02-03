package model

import "time"

type Parkn struct {
	PhoneNumber string    `bson:"phoneNumber"`
	MoveByDate  time.Time `bson:"moveByDate"`
}
