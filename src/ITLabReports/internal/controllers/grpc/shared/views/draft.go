package views

import (
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	types "github.com/RTUITLab/ITLab/proto/reports/types"
	_ "github.com/RTUITLab/ITLab/proto/reports/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func DraftFrom(draft aggregate.Draft) *types.Draft {
	return &types.Draft{
		Id:        draft.ID,
		Name:      draft.Name,
		Text:      draft.Text,
		Assignees: AssigneesFrom(draft.Assignees),
		Date:      timestamppb.New(draft.Date),
	}
}
