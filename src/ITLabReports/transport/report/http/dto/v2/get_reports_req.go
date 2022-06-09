package dto

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	queryparser "github.com/0B1t322/QueryParser"
	"github.com/0B1t322/QueryParser/typemapper"
	"github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
)

type ApprovedState string

const (
	Approved    ApprovedState = "approved"
	NotApproved ApprovedState = "notApproved"
)

type GetReportsQuery struct {
	Params *report.GetReportsParams

	ApprovedState ApprovedState
}

func (g *GetReportsQuery) NewParser() *queryparser.Parser {
	return queryparser.New(
		typemapper.NewQueryTypeFactory(),
		g.NewParseSchema(),
	)
}

func (g *GetReportsQuery) NewParseSchema() queryparser.ParseSchema {
	return queryparser.ParseSchema{
		"offset":        g.OffsetParseSchema(),
		"limit":         g.LimitParseSchema(),
		"dateBegin":     g.DateParseSchema(dateBegin),
		"dateEnd":       g.DateParseSchema(dateEnd),
		"match":         g.MatchParseSchema(),
		"sortBy":        g.SortByParseSchema(),
		"approvedState": g.ApprovedStateParseSchema(),
	}
}

type dateType string

const (
	dateBegin dateType = "dateBegin"
	dateEnd   dateType = "dateEnd"
)

func (g *GetReportsQuery) ApprovedStateParseSchema() queryparser.ParseSchemaItem {
	return queryparser.ParseSchemaItem{
		TypeMapFunc: func(field string, values []string) (interface{}, error) {
			if len(values) <= 0 {
				return nil, nil
			}
			value := values[0]
			g.ApprovedState = ApprovedState(value)

			return nil, nil
		},
	}
}
func (d dateType) GetOperation() filter.FilterOperation {
	if d == dateBegin {
		return filter.GTE
	}

	return filter.LTE
}

func (g *GetReportsQuery) DateParseSchema(dateType dateType) queryparser.ParseSchemaItem {
	return queryparser.ParseSchemaItem{
		IsRegex: false,
		ValidateValuesFunc: func(values []string) error {
			if len(values) > 0 {
				return nil
			}
			return fmt.Errorf("date not set")
		},
		TypeMapFunc: func(field string, values []string) (interface{}, error) {
			dateStr := values[0]
			date, err := time.Parse(time.RFC3339Nano, dateStr)
			if err == nil {
				g.Params.Filter.And = append(
					g.Params.Filter.And,
					&report.GetReportsFilterFieldsWithOrAnd{
						GetReportsFilterFields: report.GetReportsFilterFields{
							Date: &filter.FilterField[string]{
								Operation: dateType.GetOperation(),
								Value:     date.UTC().Format(time.RFC3339Nano),
							},
						},
					},
				)
			}

			return nil, nil
		},
	}
}

func (g *GetReportsQuery) OffsetParseSchema() queryparser.ParseSchemaItem {
	return queryparser.ParseSchemaItem{
		IsRegex: false,
		ValidateValuesFunc: func(values []string) error {
			if len(values) > 0 {
				return nil
			}
			return fmt.Errorf("offset not set")
		},
		TypeMapFunc: func(field string, values []string) (interface{}, error) {
			strOffset := values[0]

			offset, err := strconv.ParseInt(strOffset, 10, 64)
			if err == nil && offset >= 0 {
				g.Params.Offset.SetValue(offset)
			}

			return nil, nil
		},
	}
}

func (g *GetReportsQuery) LimitParseSchema() queryparser.ParseSchemaItem {
	return queryparser.ParseSchemaItem{
		IsRegex: false,
		ValidateValuesFunc: func(values []string) error {
			if len(values) > 0 {
				return nil
			}
			return fmt.Errorf("limit not set")
		},
		TypeMapFunc: func(field string, values []string) (interface{}, error) {
			strLimit := values[0]

			limit, err := strconv.ParseInt(strLimit, 10, 64)
			if err == nil && limit >= 1 {
				g.Params.Limit.SetValue(limit)
			}

			return nil, nil
		},
	}
}

// MatchParseSchema Describe match format field_1:value_1
func (g *GetReportsQuery) MatchParseSchema() queryparser.ParseSchemaItem {
	return queryparser.ParseSchemaItem{
		IsRegex: false,
		TypeMapFunc: func(field string, values []string) (interface{}, error) {
			for _, matchValue := range values {
				fieldValue := strings.SplitN(matchValue, ":", 2)
				if len(fieldValue) != 2 {
					continue
				}
				field := fieldValue[0]
				value := fieldValue[1]

				g.SetMatchField(field, value)
			}

			return nil, nil
		},
	}
}

// SortByParseSchema Describe match format field_1:order where order = asc|desc
func (g *GetReportsQuery) SortByParseSchema() queryparser.ParseSchemaItem {
	return queryparser.ParseSchemaItem{
		IsRegex: false,
		TypeMapFunc: func(field string, values []string) (interface{}, error) {
			for _, sortField := range values {
				fieldValue := strings.Split(sortField, ":")
				if len(fieldValue) != 2 {
					continue
				}
				field := fieldValue[0]
				value := ordertype.OrderTypeFromString(fieldValue[1])

				g.SetSortField(field, value)
			}

			return nil, nil
		},
	}
}

func (g *GetReportsQuery) SetSortField(field string, order ordertype.OrderType) {
	switch field {
	case "date":
		g.Params.Filter.DateSort.SetValue(order)
	case "name":
		g.Params.Filter.NameSort.SetValue(order)
	}
}

func (g *GetReportsQuery) SetMatchField(field string, value string) {
	switch field {
	case "name":
		g.Params.Filter.Name = &filter.FilterField[string]{
			Operation: filter.LIKE,
			Value:     value,
		}
	case "text":
		// TODO if needed
	case "date":
		date, err := time.Parse(time.RFC3339Nano, value)
		if err != nil {
			return
		}

		g.Params.Filter.Date = &filter.FilterField[string]{
			Operation: filter.EQ,
			Value:     date.UTC().Format(time.RFC3339Nano),
		}
	case "assignees.implementer":
		g.Params.Filter.Implementer = &filter.FilterField[string]{
			Operation: filter.EQ,
			Value:     value,
		}
	case "assignees.reporter":
		g.Params.Filter.Reporter = &filter.FilterField[string]{
			Operation: filter.EQ,
			Value:     value,
		}
	}
}

type GetReportsReq struct {
	Query GetReportsQuery
}

func (g *GetReportsReq) IsOnlyApprovedReports() bool {
	return g.Query.ApprovedState == Approved
}

func (g *GetReportsReq) SetOnlyApprovedReports(ids ...string) {
	g.Query.Params.Filter.ReportsId = &filter.FilterField[[]string]{
		Value:     ids,
		Operation: filter.IN,
	}
}

func (g *GetReportsReq) IsOnlyNotApprovedReports() bool {
	return g.Query.ApprovedState == NotApproved
}

func (g *GetReportsReq) SetOnlyNotApprovedReports(ids ...string) {
	g.Query.Params.Filter.ReportsId = &filter.FilterField[[]string]{
		Value:     ids,
		Operation: filter.NIN,
	}
}

func (g *GetReportsReq) SetImplementerAndReporter(implementer, reporter string) {
	// If this method call it's mean user don't have access to another reports
	// So nill filters
	g.Query.Params.Filter.Implementer = nil
	g.Query.Params.Filter.Reporter = nil

	g.Query.Params.Filter.Or = append(
		g.Query.Params.Filter.Or,
		&report.GetReportsFilterFieldsWithOrAnd{
			GetReportsFilterFields: report.GetReportsFilterFields{
				Reporter: &filter.FilterField[string]{
					Value:     reporter,
					Operation: filter.EQ,
				},
			},
		},
		&report.GetReportsFilterFieldsWithOrAnd{
			GetReportsFilterFields: report.GetReportsFilterFields{
				Reporter: &filter.FilterField[string]{
					Value:     implementer,
					Operation: filter.EQ,
				},
			},
		},
	)
}

func (g *GetReportsReq) ToEndpointReq() *reqresp.GetReportsReq {
	req := &reqresp.GetReportsReq{
		Params: g.Query.Params,
	}

	return req
}

func DecodeGetReportsReq(
	ctx context.Context,
	r *http.Request,
) (*GetReportsReq, error) {
	req := &GetReportsReq{
		Query: GetReportsQuery{
			Params: &report.GetReportsParams{
				Filter: &report.GetReportsFilter{},
			},
		},
	}

	req.Query.NewParser().ParseUrlValues(r.URL.Query())

	return req, nil
}
