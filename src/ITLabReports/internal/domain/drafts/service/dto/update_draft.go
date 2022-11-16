package dto

import (
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/samber/mo"
)

type UpdateDraftReq struct {
	ID          string
	Name        mo.Option[string]
	Text        mo.Option[string]
	Implementer mo.Option[string]
	By          aggregate.User
}
