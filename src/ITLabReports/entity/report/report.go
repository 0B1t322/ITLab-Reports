package report

import (
	"time"
)

type ReportState string

const (
	ReportStateCreated ReportState = "created"
	ReportStatePaid    ReportState = "paid"
)

type Report struct {
	ID    string
	Name  string
	Date  time.Time
	Text  string
	State ReportState
}
