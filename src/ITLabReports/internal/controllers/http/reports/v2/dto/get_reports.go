package dto

import (
	"time"

	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/match"
	"github.com/RTUITLab/ITLab-Reports/internal/controllers/http/shared/view"
)

type ReportMatchFields string

const (
	ReportMatchFields_Name        ReportMatchFields = "name"
	ReportMatchFields_Date        ReportMatchFields = "date"
	ReportMatchFields_Reporter    ReportMatchFields = "assignees.reporter"
	ReportMatchFields_Implementer ReportMatchFields = "assignees.implementer"
)

var (
	Name_MatchParam = match.NewMatchDesc(
		ReportMatchFields_Name,
		func(s string) (string, error) {
			return s, nil
		},
	)

	Date_MatchParam = match.NewMatchDesc(
		ReportMatchFields_Date,
		func(s string) (time.Time, error) {
			return time.Parse(time.RFC3339, s)
		},
	)

	Reporter_MatchParam = match.NewMatchDesc(
		ReportMatchFields_Reporter,
		func(s string) (string, error) {
			return s, nil
		},
	)

	Implementer_MatchParam = match.NewMatchDesc(
		ReportMatchFields_Implementer,
		func(s string) (string, error) {
			return s, nil
		},
	)
)

type ReportSortFields string

const (
	ReportSortFields_Name ReportSortFields = "name"
	ReportSortFields_Date ReportSortFields = "date"
)

type ApprovedState string

const (
	ApprovedState_All      ApprovedState = ""
	ApprovedState_Approved ApprovedState = "approved"
	ApprovedState_Not      ApprovedState = "notApproved"
)

type (
	GetReportsReq struct {
		DateBegin     time.Time     `json:"dateBegin"     validate:"optional" form:"dateBegin"`
		DateEnd       time.Time     `json:"dateEnd"       validate:"optional" form:"dateEnd"`
		Offset        int64         `json:"offset"        validate:"optional" form:"offset"`
		Limit         int64         `json:"limit"         validate:"optional" form:"limit"`
		Match         []string      `json:"match"         validate:"optional" form:"match"         swaggerignore:"true"`
		SortBy        []string      `json:"sortBy"        validate:"optional" form:"sortBy"        swaggerignore:"true"`
		ApprovedState ApprovedState `json:"approvedState" validate:"optional" form:"approvedState"`
	}

	GetReportsResp struct {
		Count       int64             `json:"count"`
		Reports     []view.ReportView `json:"items"`
		HasMore     bool              `json:"hasMore"`
		Limit       int64             `json:"limit"`
		Offset      int64             `json:"offset"`
		TotalResult int64             `json:"totalResult"`
	}
)
