package dto

type UpdateDraftReq struct {
	ID          string `json:"-"           uri:"id" swaggerignore:"true"`
	Name        string `json:"name"                                      validate:"optional"`
	Text        string `json:"text"                                      validate:"optional"`
	Implementer string `json:"implementer"                               validate:"optional"`
}
