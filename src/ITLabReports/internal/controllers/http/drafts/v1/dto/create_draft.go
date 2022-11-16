package dto

type CreateDraftReq struct {
	Name        string `json:"name" validate:"required"`
	Text        string `json:"text" validate:"required"`
	Implementer string `json:"-"                        form:"implementer" swaggerignore:"true"`
}
