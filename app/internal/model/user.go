package model

import "time"

type User struct {
	Id                int
	Username          string
	Password          string
	CreateTime        time.Time
	PersonalSignature string
}
