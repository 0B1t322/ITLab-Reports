package models

import (
	"time"

	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/RTUITLab/ITLab-Reports/internal/models/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DraftFields string

func (d DraftFields) String() string {
	return string(d)
}

const (
	DraftFieldsID        DraftFields = "_id"
	DraftFieldsName      DraftFields = "name"
	DraftFieldsDate      DraftFields = "date"
	DraftFieldsText      DraftFields = "text"
	DraftFieldsAssignees DraftFields = "assignees"
)

type Draft struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Date      time.Time          `bson:"date"`
	Text      string             `bson:"text"`
	Assignees Assignees          `bson:"assignees"`
}

func NewDraftModel(draft aggregate.Draft) (Draft, error) {
	id, err := primitive.ObjectIDFromHex(draft.ID)
	if err != nil {
		return Draft{}, err
	}

	return Draft{
		ID:        id,
		Name:      draft.Name,
		Date:      draft.Date.UTC(),
		Text:      draft.Text,
		Assignees: NewAssigneesModel(draft.Assignees),
	}, nil
}

func DraftFromModel(draft Draft) aggregate.Draft {
	return aggregate.Draft{
		Report: &entity.Report{
			ID:   draft.ID.Hex(),
			Name: draft.Name,
			Date: draft.Date,
			Text: draft.Text,
		},
		Assignees: AssigneesFromModel(draft.Assignees),
	}
}
