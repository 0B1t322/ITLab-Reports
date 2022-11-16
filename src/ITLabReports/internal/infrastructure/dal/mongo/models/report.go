package models

import (
	"time"

	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/RTUITLab/ITLab-Reports/internal/models/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportState string

func (r ReportState) String() string {
	return string(r)
}

const (
	ReportStateCreated ReportState = "created"
	ReportStatePaid    ReportState = "paid"
)

type ReportFields string

func (r ReportFields) String() string {
	return string(r)
}

const (
	ReportFieldsID        ReportFields = "_id"
	ReportFieldsName      ReportFields = "name"
	ReportFieldsDate      ReportFields = "date"
	ReportFieldsText      ReportFields = "text"
	ReportFieldsState     ReportFields = "state"
	ReportFieldsAssignees ReportFields = "assignees"
)

type Report struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Date      time.Time          `bson:"date"`
	Text      string             `bson:"text"`
	Assignees Assignees          `bson:"assignees"`
	State     ReportState        `bson:"state"`
}

func NewReportModel(report aggregate.Report) (Report, error) {
	id, err := primitive.ObjectIDFromHex(report.ID)
	if err != nil {
		return Report{}, err
	}

	return Report{
		ID:        id,
		Name:      report.Name,
		Date:      report.Date.UTC(),
		Text:      report.Text,
		Assignees: NewAssigneesModel(report.Assignees),
		State:     ReportState(report.State),
	}, nil
}

func ReportFromModel(report Report) aggregate.Report {
	return aggregate.Report{
		Report: &entity.Report{
			ID:   report.ID.Hex(),
			Name: report.Name,
			Date: report.Date.UTC(),
			Text: report.Text,
		},
		Assignees: AssigneesFromModel(report.Assignees),
		State:     aggregate.ReportState(report.State),
	}
}
