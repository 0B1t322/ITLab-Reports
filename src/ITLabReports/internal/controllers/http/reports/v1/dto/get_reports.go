package dto

type ReportsSortedByFields string

const (
	ReportsSortedByFields_Name ReportsSortedByFields = "name"
	ReportsSortedByFields_Date ReportsSortedByFields = "date"
)

type GetReportsReq struct {
	SortedBy ReportsSortedByFields `json:"sorted_by" form:"sorted_by"`
}
