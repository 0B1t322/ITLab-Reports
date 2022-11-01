package dto

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	queryparser "github.com/0B1t322/QueryParser"
	"github.com/0B1t322/QueryParser/typemapper"
	"github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
	"github.com/RTUITLab/ITLab-Reports/transport/report/reqresp"
	"github.com/samber/mo"
)

type GetDraftQuery struct {
	Params *report.GetReportsParams
}

func (g *GetDraftQuery) NewParser() *queryparser.Parser {
	return queryparser.New(
		typemapper.NewQueryTypeFactory(),
		g.NewParseSchema(),
	)
}

func (g *GetDraftQuery) NewParseSchema() queryparser.ParseSchema {
	return queryparser.ParseSchema{
		"offset": g.OffsetParseSchema(),
		"limit":  g.LimitParseSchema(),
	}
}

func (g *GetDraftQuery) OffsetParseSchema() queryparser.ParseSchemaItem {
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
				g.Params.Offset = mo.Some(offset)
			}

			return nil, nil
		},
	}
}

func (g *GetDraftQuery) LimitParseSchema() queryparser.ParseSchemaItem {
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
				g.Params.Limit = mo.Some(limit)
			}

			return nil, nil
		},
	}
}

type GetDraftsReq struct {
	Query  GetDraftQuery
	UserID string
}

func (g *GetDraftsReq) SetUserID(id string) {
	g.UserID = id
}

func (g *GetDraftsReq) ToEndpointReq() *reqresp.GetReportsReq {
	g.Query.Params.Filter.SortParams = append(
		g.Query.Params.Filter.SortParams,
		report.GetReportsSort{
			DateSort: mo.Some[ordertype.OrderType](ordertype.DESC),
		},
	)

	g.Query.Params.Filter.GetReportsFilterFieldsWithOrAnd.Reporter = &filter.FilterField[string]{
		Value:     g.UserID,
		Operation: filter.EQ,
	}

	return &reqresp.GetReportsReq{
		Params: g.Query.Params,
	}
}

func DecodeGetDraftsReq(
	ctx context.Context,
	r *http.Request,
) (*GetDraftsReq, error) {
	req := &GetDraftsReq{
		Query: GetDraftQuery{
			Params: &report.GetReportsParams{
				Filter: &report.GetReportsFilter{},
			},
		},
	}

	req.Query.NewParser().ParseUrlValues(r.URL.Query())

	return req, nil
}
