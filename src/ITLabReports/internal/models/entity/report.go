package entity

import "time"

type Report struct {
	ID   string
	Name string
	Date time.Time
	Text string
}
