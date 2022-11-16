package models

import "github.com/RTUITLab/ITLab-Reports/internal/models/valueobject"

type AssigneesFields string

func (a AssigneesFields) String() string {
	return string(a)
}

const (
	AssigneesFieldsReporter    AssigneesFields = "reporter"
	AssigneesFieldsImplementer AssigneesFields = "implementer"
)

type Assignees struct {
	Reporter    string `bson:"reporter"`
	Implementer string `bson:"implementer"`
}

func NewAssigneesModel(assignees valueobject.Assignees) Assignees {
	return Assignees{
		Reporter:    assignees.Reporter,
		Implementer: assignees.Implementer,
	}
}

func AssigneesFromModel(assignees Assignees) valueobject.Assignees {
	return valueobject.Assignees{
		Reporter:    assignees.Reporter,
		Implementer: assignees.Implementer,
	}
}
