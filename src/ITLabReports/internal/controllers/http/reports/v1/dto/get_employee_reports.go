package dto

import "time"

type GetEmployeeReportsReq struct {
	DateBegin  time.Time `json:"dateBegin" form:"dateBegin"`
	DateEnd    time.Time `json:"dateEnd"   form:"dateEnd"`
	EmployeeID string    `json:"employee"                   uri:"employee"`
}
