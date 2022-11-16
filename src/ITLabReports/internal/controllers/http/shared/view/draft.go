package view

import (
	"time"

	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
)

type DraftView struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Text      string        `json:"text"`
	Assignees AssigneesView `json:"assignees"`
	Date      time.Time     `json:"date"`
}

type DraftsView struct {
	Drafts []DraftView `json:"drafts"`
}

func DraftViewFrom(from aggregate.Draft) DraftView {
	return DraftView{
		ID:        from.ID,
		Name:      from.Name,
		Text:      from.Text,
		Assignees: AssigneesViewFrom(from.Assignees),
		Date:      from.Date.UTC(),
	}
}

func DraftsViewFrom(from []aggregate.Draft) DraftsView {
	var drafts []DraftView

	for _, draft := range from {
		drafts = append(drafts, DraftViewFrom(draft))
	}

	return DraftsView{
		Drafts: drafts,
	}
}
