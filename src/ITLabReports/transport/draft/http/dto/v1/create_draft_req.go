package dto

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"

	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	entassignees "github.com/RTUITLab/ITLab-Reports/entity/assignees"
	entreport "github.com/RTUITLab/ITLab-Reports/entity/report"
)

type CreateDraftReq struct {
	Name string	`json:"name,omitempty"`
	Text string `json:"text,omitempty"`
	Implementor string `json:"-" swaggeringore:"true"`
	Reporter string `json:"-" swaggeringore:"true"`
}

func (c *CreateDraftReq) SetReporter(r string) {
	c.Reporter = r
}

func (c *CreateDraftReq) SetImplementor(i string) {
	c.Implementor = i
}

func (c *CreateDraftReq) GetImplementor() string {
	return c.Implementor
}

func (c *CreateDraftReq) ToEndpointReq() *reqresp.CreateReportReq {
	if c.Implementor == "" {
		c.Implementor = c.Reporter
	}
	
	rep := &report.Report{
		Report: &entreport.Report{
			Name: c.Name,
			Text: c.Text,
			Date: time.Now().UTC().Round(time.Millisecond),
		},
		Assignees: &entassignees.Assignees{
			Reporter: c.Reporter,
			Implementer: c.Implementor,
		},
	}

	return &reqresp.CreateReportReq{
		Report: rep,
	}
}

func DecodeCreateDraftReq(
	ctx context.Context,
	r *http.Request,
) (*CreateDraftReq, error) {
	values := r.URL.Query()

	var (
		implementer = values.Get("implementer")
	)

	req := &CreateDraftReq{
		Implementor: implementer,
	}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	return req, nil
}