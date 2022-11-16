package view

import "github.com/RTUITLab/ITLab-Reports/internal/models/valueobject"

type AssigneesView struct {
	Reporter    string `json:"reporter"`
	Implementer string `json:"implementer"`
}

func AssigneesViewFrom(from valueobject.Assignees) AssigneesView {
	return AssigneesView{
		Reporter:    from.Reporter,
		Implementer: from.Implementer,
	}
}
