package dto

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	entassignees "github.com/RTUITLab/ITLab-Reports/entity/assignees"
	entreport "github.com/RTUITLab/ITLab-Reports/entity/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/clarketm/json"
	"github.com/gorilla/mux"
)

type CreateReportReq struct {
	Name *string `json:"name,omitempty" swaggertype:"string" extensions:"x-nullable`
	Text string `json:"text"`
	Implementor string `json:"-" swaggeringore:"true"`
	Reporter string `json:"-" swaggeringore:"true"`
}

func (c *CreateReportReq) SetReporter(r string) {
	c.Reporter = r
}

func (c *CreateReportReq) GetName() string {
	if c.Name != nil {
		return *c.Name
	} else {
		splitedText := strings.Split(c.Text, "@\n\t\n@")
		if len(splitedText) < 2 {
			return ""
		}

		return splitedText[0]
	}
}

func (c *CreateReportReq) GetText() string {
	if c.Name != nil {
		return c.Text
	} else {
		splitedText := strings.Split(c.Text, "@\n\t\n@")
		if len(splitedText) < 2 {
			return ""
		}

		return splitedText[1]
	}
}

func (c *CreateReportReq) ToEndpointReq() *reqresp.CreateReportReq {
	rep := &report.Report{
		Report: &entreport.Report{
			Name: c.GetName(),
			Date: time.Now().UTC(),
			Text: c.GetText(),
		},
		Assignees: &entassignees.Assignees{
			Implementer: c.Implementor,
			Reporter: c.Reporter,
		},
	}

	return &reqresp.CreateReportReq{
		Report: rep,
	}
}

func DecodeCreateReportReq(
	ctx context.Context,
	r *http.Request,
) (*CreateReportReq, error) {
	vars := mux.Vars(r)

	var (
		implementer = vars["implementor"]
	)

	req := &CreateReportReq{
		Implementor: implementer,
	}


	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}


	return req, nil
}