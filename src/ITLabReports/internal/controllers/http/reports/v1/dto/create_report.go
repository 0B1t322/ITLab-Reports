package dto

type CreateReportReq struct {
	Name        string  `json:"name"`
	Text        string  `json:"text"`
	Implementer *string `json:"-"    form:"implementer" swaggerignore:"true"`
}
