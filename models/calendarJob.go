package models

import "time"

type CalendarJob struct {
	Id       int
	Bth      bool
	SubEmail string
	SubName  string
	Text     string
	Date     time.Time
}
