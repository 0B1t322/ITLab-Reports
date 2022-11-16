package dto

import (
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/samber/mo"
)

type CreateDraftReq struct {
	Name        string
	Text        string
	Implementer mo.Option[string]
	By          aggregate.User
}
