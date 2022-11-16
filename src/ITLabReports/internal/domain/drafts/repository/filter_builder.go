/*
Provide struct and methods to build filter query
Code generatated with repogen
Do not Edit
*/
package drafts

import (
	"github.com/0B1t322/RepoGen/pkg/filter"
	"github.com/samber/mo"
)

type filterQuery struct {
	FilterQuery
}

func Query() filterQuery {
	return filterQuery{}
}

func (q filterQuery) Build() FilterQuery {
	return q.FilterQuery
}

func (q filterQuery) Or(es ...expression) filterQuery {
	for _, e := range es {
		q.FilterQuery.Or = append(q.FilterQuery.Or, Query().Expression(e).FilterQuery)
	}
	return q
}

func (q filterQuery) And(es ...expression) filterQuery {
	for _, e := range es {
		q.FilterQuery.And = append(q.FilterQuery.And, Query().Expression(e).FilterQuery)
	}
	return q
}

func (q filterQuery) Expression(e expression) filterQuery {
	q.FilterQuery.Expression = e.FilterFields
	return q
}

type expression struct {
	FilterFields
}

func Expression() expression {
	return expression{}
}

func (e expression) Build() FilterFields {
	return e.FilterFields
}

func (e expression) ID(id string, op filter.FilterOperation) expression {
	e.FilterFields.ID = mo.Some(filter.New(id, op))
	return e
}

func (e expression) Implementer(implementer string, op filter.FilterOperation) expression {
	e.FilterFields.Implementer = mo.Some(filter.New(implementer, op))
	return e
}

func (e expression) Reporter(reporter string, op filter.FilterOperation) expression {
	e.FilterFields.Reporter = mo.Some(filter.New(reporter, op))
	return e
}
