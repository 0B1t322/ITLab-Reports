package mongo

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Builder struct {
	raw 	bson.M
	op		opearation
}

func NewBuilder(
	op	opearation,
) Builder {
	return Builder{
		raw: bson.M{},
		op: op,
	}
}

// Predicate is filter predicate
type Predicate struct {
	raw		bson.M
}

func (p Predicate) MarshalBSON() ([]byte, error) {
	return bson.Marshal(p.raw)
}

func (p *Predicate) BSON() bson.M {
	return p.raw
}

func (p *Predicate) Object() interface{} {
	return p.BSON()
}

func (p *Predicate) append(field string, op string, value interface{}) *Predicate {
	if _, find := p.raw[field]; find {
		switch val := p.raw[field].(type) {
		case bson.M:
			val[op] = value
		}
	} else {
		p.raw[field] = bson.M{op: value}
	}
	return p
}

const (
	eq 					= "$eq"
	ne					= "$ne"
	or 					= "$or"
	and 				= "$and"
	not					= "$not"
	gt					= "$gt"
	gte					= "$gte"
	lt					= "$lt"
	lte					= "$lte"
	in					= "$in"
	nin					= "$nin"
	exists				= "$exists"
	regex				= "$regex"
	elemMatch			= "$elemMatch"
)

func P() *Predicate {
	return &Predicate{raw: bson.M{}}
}

// return {"$or": [{...},...,{...}]}
func Or(preds ...*Predicate) *Predicate {
	p := P()
	array := bson.A{}
	for _, pred := range preds {
		if pred == nil {
			continue
		}
		array = append(array, pred.raw)
	}
	p.raw[or] = array
	return p
}

func (p *Predicate) Or(preds ...*Predicate) *Predicate {
	if _, find := p.raw[or]; find {
		array := p.raw[or].(bson.A)
		for _, pred := range preds {
			if pred == nil {
				continue
			}
			array = append(array, pred.raw)
		}
		p.raw[or] = array
		return p
	}
	p.raw[or] = Or(preds...).raw[or]

	return p
}

// EQ returns {"field": value}
func EQ(field string, value interface{}) *Predicate {
	return P().EQ(field, value)
}

// AEQ return {"$eq": ["field", value]}
func (p *Predicate) AEQ(field string, value interface{}) *Predicate {
	p.raw[eq] = Array().AddElems(field, value).Array()
	return p
}

// AEQ return {"$eq": ["field", value]}
func AEQ(field string, value interface{}) *Predicate {
	return P().AEQ(field, value)
}

// NEQ returns {"field": {"$ne": value}}
func NEQ(field string, value interface{}) *Predicate {
	return P().NEQ(field, value)
}

// ANEQ return {"$neq": ["field", value]}
func (p *Predicate) ANEQ(field string, value interface{}) *Predicate {
	p.raw[ne] = Array().AddElems(field, value).Array()
	return p
}

// ANEQ return {"$neq": ["field", value]}
func ANEQ(field string, value interface{}) *Predicate {
	return P().ANEQ(field, value)
}

// EQOper return {"field": {"$eq": value}}
func EQOper(field string, value interface{}) *Predicate {
	return P().EQOper(field, value)
}

// EQ returns {"field": value}
func (p *Predicate) EQ(field string, value interface{}) *Predicate {
	p.raw[field] = value
	return p
}

// EQOper return {"field": {"$eq": value}}
func (p *Predicate) EQOper(field string, value interface{}) *Predicate {
	return p.append(field, eq, value)
}

// NEQ returns {"field": {"$ne": value}}
func (p *Predicate) NEQ(field string, value interface{}) *Predicate {
	return p.append(field, ne, value)
}

// LT returns {"field": {"$lt": value}}
func LT(field string, value interface{}) *Predicate {
	return P().LT(field, value)
}

// ALT return {"$ne": ["field", value]}
func (p *Predicate) ALT(field string, value interface{}) *Predicate {
	p.raw[lt] = Array().AddElems(field, value).Array()
	return p
}

// ALT return {"$ne": ["field", value]}
func ALT(field string, value interface{}) *Predicate {
	return P().ALT(field, value)
}

// LT returns {"field": {"$lt": value}}
func (p *Predicate) LT(field string, value interface{}) *Predicate {
	return p.append(field, lt, value)
}

// LTE returns {"field": {"$lte": value}}
func LTE(field string, value interface{}) *Predicate {
	return P().LTE(field, value)
}

// LTE returns {"field": {"$lte": value}}
func (p *Predicate) LTE(field string, value interface{}) *Predicate {
	return p.append(field, lte, value)
}

// ALTE return {"$lte": ["field", value]}
func (p *Predicate) ALTE(field string, value interface{}) *Predicate {
	p.raw[lte] = Array().AddElems(field, value).Array()
	return p
}

func ALTE(field string, value interface{}) *Predicate {
	return P().ALTE(field, value)
}

// GTE returns {"field": {"$gte": value}}
func GTE(field string, value interface{}) *Predicate {
	return P().GTE(field, value)
}

// GTE returns {"field": {"$gte": value}}
func (p *Predicate) GTE(field string, value interface{}) *Predicate {
	return p.append(field, gte, value)
}

// AGTE return {"$gte": ["field", value]}
func (p *Predicate) AGTE(field string, value interface{}) *Predicate {
	p.raw[gte] = Array().AddElems(field, value).Array()
	return p
}

// AGTE return {"$gte": ["field", value]}
func AGTE(field string, value interface{}) *Predicate {
	return P().AGTE(field, value)
}

// GT returns {"field": {"$gt": value}}
func GT(field string, value interface{}) *Predicate {
	return P().GT(field, value)
}

// AGT return {"$gt": ["field", value]}
func (p *Predicate) AGT(field string, value interface{}) *Predicate {
	p.raw[gt] = Array().AddElems(field, value).Array()
	return p
}

// AGT return {"$gt": ["field", value]}
func AGT(field string, value interface{}) *Predicate {
	return P().AGT(field, value)
}

// GT returns {"field": {"$gt": value}}
func (p *Predicate) GT(field string, value interface{}) *Predicate {
	return p.append(field, gt, value)
}

// In returns {"field": {"$in": [...]}}
func In(field string, values ...interface{}) *Predicate {
	return P().In(field, values...)
}

// In returns {"field": {"$in": [...]}}
func (p *Predicate) In(field string, values ...interface{}) *Predicate {
	if len(values) == 0 {
		return p
	}
	array := bson.A{}
	array = append(array, values...)

	if raw, ok := p.raw[field].(bson.M); ok {
		value, find := raw[in]
		if find {
			prevaArray := value.(bson.A)
			prevaArray = append(prevaArray, array...)
			array = prevaArray
		}
	}

	return p.append(field, in, array)
}

// NotIn returns {"field": {"$nin": [...]}}
func NotIn(field string, values ...interface{}) *Predicate {
	return P().NotIn(field, values...)
}

// NotIn returns {"field": {"$nin": [...]}}
func (p *Predicate) NotIn(field string, values ...interface{}) *Predicate {
	if len(values) == 0 {
		return p
	}
	array := bson.A{}
	array = append(array, values...)
	if raw, ok := p.raw[field].(bson.M); ok {
		value, find := raw[nin]
		if find {
			prevaArray := value.(bson.A)
			prevaArray = append(prevaArray, array...)
			array = prevaArray
		}
	}

	return p.append(field, nin, array)
}

// Exist return {"field": {"$exists": true}}
func Exist(field string) *Predicate {
	return P().Exist(field)
}

// Exist return {"field": {"$exists": true}}
func (p *Predicate) Exist(field string) *Predicate {
	return p.append(field, exists, true)
}

// NotExist return {"field": {"$exists": false}}
func NotExist(field string) *Predicate {
	return P().Exist(field)
}

// NotExist return {"field": {"$exists": false}}
func (p *Predicate) NotExist(field string) *Predicate {
	return p.append(field, exists, false)
}

func ElemMatch(field string, pred *Predicate) *Predicate {
	return P().ElemMatch(field, pred)
}

func (p *Predicate) ElemMatch(field string, pred *Predicate) *Predicate {
	p.append(field, elemMatch, pred.raw)
	return p
}

type RegexOptions int64
const (
	I = 1<<iota
	M
	X
	S
)
func option(opt RegexOptions) string {
	options := strings.Builder{}
	if opt & I == I {
		options.WriteString("i")
	}
	if opt & M == M {
		options.WriteString("m")
	}
	if opt & X == X {
		options.WriteString("x")
	}
	if opt & S == S {
		options.WriteString("s")
	}
	return options.String()
}

/* 
Like return 
	{"field": {"$regex": pattern}}

if opts not 0 return 
	{"field": {"$regex": pattern, "$options": "..."}}
 
If you want to add more than option use 

	Like("some_field", "a", I | M)
*/
func (p *Predicate) Like(field string, pattern string, opts RegexOptions) *Predicate {
	new := p.append(field, regex, pattern)
	if opts != 0 {
		new = new.append(field, "$options", option(opts))
	}
	return new
}

/* 
Like return 
	{"field": {"$regex": pattern}}

if opts not 0 return 
	{"field": {"$regex": pattern, "$options": "..."}}
 
If you want to add more than option use 

	Like("some_field", "a", I | M)
*/
func Like(field string, pattern string, opts RegexOptions) *Predicate {
	return P().Like(field, pattern, opts)
}

// return {"$not": {...}}
func Not(pred *Predicate) *Predicate {
	return P().Not(pred)
}

// return {"$not": {...}}
func (p *Predicate) Not(pred *Predicate) *Predicate {
	p.raw[not] = pred.raw
	return p
}

// return {"$and": [{...},...,{...}]}
func And(preds ...*Predicate) *Predicate {
	p := P()
	array := bson.A{}
	for _, pred := range preds {
		if pred == nil {
			continue
		}
		array = append(array, pred.raw)
	}
	p.raw[and] = array
	return p
}

func (p *Predicate) And(preds ...*Predicate) *Predicate {
	if _, find := p.raw[and]; find {
		array := p.raw[and].(primitive.A)
		for _, pred := range preds {
			if pred == nil {
				continue
			}

			array = append(array, pred.raw)
		}
		p.raw[and] = array
		return p
	}
	p.raw[and] = And(preds...).raw[and]
	return p
}

type UpdateBuilder struct {
	Builder
	filter		*Predicate
	update		interface{}
	collection	string
}

// Return update builder
// in standart operation set to updateMany
func Update(collection string) *UpdateBuilder { 
	return &UpdateBuilder{
		Builder: NewBuilder(UpdateMany),
	}
}

// Set filter to update operation
// if filter exist wrap both filter by and
func (u *UpdateBuilder) Filter(p *Predicate) *UpdateBuilder {
	if u.filter != nil {
		u.filter = And(u.filter, p)
	} else {
		u.filter = p
	}
	return u
}

// Update only one matched field
func (u *UpdateBuilder) One() *UpdateBuilder {
	u.op = UpdateOne
	return u
}

func (u *UpdateBuilder) Set(
	s func(s *SetBuidler),
) *UpdateBuilder {
	setBuidler := Set()
	s(setBuidler)
	u.update = setBuidler.raw
	return u
}

type SetBuidler struct {
	raw bson.M
}

func (s SetBuidler) MarshalBSON() ([]byte, error) {
	return bson.Marshal(s.raw)
}

func Set() *SetBuidler {
	return &SetBuidler{
		raw: bson.M{"$set": bson.M{}},
	}
}

func (s *SetBuidler) getSet() interface{} {
	return s.raw["$set"]
}

type operationBuilder struct {
	operation string
	raw bson.M
}

// Create new opeataion builder
// Use to create builders for operations for example $sum
func newOpearionBuilder(operation string) *operationBuilder {
	return &operationBuilder{
		raw: bson.M{fmt.Sprintf("$%s", operation): bson.M{}},
		operation: operation,
	}
}

func (o *operationBuilder) getOpearation() interface{} {
	return o.raw[o.operation]
}

// add to "$operation" key-value
// if this field exist ovetwrite value
// if you want to add field as object user ObjectBuilder
// return:
// 	{"$operation": {"field": value}}
func (o *operationBuilder) SetField(field string, value interface{}) *operationBuilder {
	o.getOpearation().(bson.M)[field] = value
	return o
}

// add to "$operation" key-value in slice field
// if this field exist ovetwrite value
// if you want to add field as object user ObjectBuilder
// return:
// 	{"$operation": {"sliceField.$.field": value}}
func (o *operationBuilder) SetSliceField(
	sliceField	string,
	field 		string, 
	value 		interface{},
) *operationBuilder {
	o.getOpearation().(bson.M)[fmt.Sprintf("%s.$.%s",sliceField, field)] = value
	return o
}

// add to "$operation" key-value in slice like singe value
// if this field exist ovetwrite value
// if you want to add field as object user ObjectBuilder
// return:
// 	{"$operation": {"sliceField.$": value}}
func (s *operationBuilder) SetSliceValue(
	sliceField	string,
	value		interface{},
) *operationBuilder {
	s.getOpearation().(bson.M)[fmt.Sprintf("%s.$", sliceField)] = value
	return s
}

type ObjectBuilder struct {
	raw bson.M
}

// return new object builder
func newObjectBuilder() *ObjectBuilder {
	return &ObjectBuilder{raw: bson.M{}}
}

func Object() *ObjectBuilder {
	return newObjectBuilder()
}

type toObject interface {
	Object() interface{}
}

type toArray interface {
	Array() interface{}
}

// set field to object
// if you want to use to add in field new object use it like:
//	o.AddField("field": NewObjectBuidler().AddField("new_field": value).Object())
// return
// 	{"field": value}
func (o *ObjectBuilder) AddField(field string, value interface{}) *ObjectBuilder {
	if objBuilder, ok := value.(toObject); ok {
		o.raw[field] = objBuilder.Object()
	} else if array, ok := value.(toArray); ok {
		o.raw[field] = array.Array()
	} else {
		o.raw[field] = value
	}
	return o
}

func (o *ObjectBuilder) AddFieldIf(field string, value interface{}, cond func() bool) *ObjectBuilder {
	if cond() {
		o.AddField(field, value)
	}
	return o
}

// set slice field
// return
// 	{"field": [...]}
func (o *ObjectBuilder) AddSliceField(field string, values ...interface{}) *ObjectBuilder {
	slice := bson.A{}
	for _, value := range values {
		if objBuidler, ok := value.(toObject); ok {
			slice = append(slice, objBuidler.Object())
		} else {
			slice = append(slice, value)
		}
	}
	o.raw[field] = slice
	return o
}

func (o *ObjectBuilder) Object() interface{} {
	return o.raw
}

type ArrayBuidler struct {
	raw bson.A
}

func newArrayBuilder() *ArrayBuidler {
	return &ArrayBuidler{raw: bson.A{}}
}

func Array() *ArrayBuidler {
	return newArrayBuilder()
}

func (a *ArrayBuidler) AddElem(elem interface{}) *ArrayBuidler {
	if obj, ok := elem.(toObject); ok {
		a.raw = append(a.raw, obj.Object())
	} else {
		a.raw = append(a.raw, elem)
	}
	return a
}

func (a *ArrayBuidler) AddElems(elems ...interface{}) *ArrayBuidler {
	for _, elem := range elems {
		a.AddElem(elem)
	}
	return a
}

func (a *ArrayBuidler) Array() interface{} {
	return a.raw
}

// add to "$set" key-value
// if this field exist ovetwrite value
// if you want to add field as object user ObjectBuilder
// return:
// 	{"$set": {"field": value}}
func (s *SetBuidler) SetField(field string, value interface{}) *SetBuidler {
	s.getSet().(bson.M)[field] = value
	return s
}

// add to "$set" key-value if condition func return true
// 
// if this field exist ovetwrite value
// 
// if you want to add field as object user ObjectBuilder
// 
// return:
// 	{"$set": {"field": value}}
func (s *SetBuidler) SetFieldIf(field string, value interface{}, cond func() bool) *SetBuidler {
	if !cond() {
		return s
	}
	return s.SetField(field, value)
}

// add to "$set" key-value in slice field
// if this field exist ovetwrite value
// if you want to add field as object user ObjectBuilder
// return:
// 	{"$set": {"sliceField.$.field": value}}
func (s *SetBuidler) SetSliceField(
	sliceField	string,
	field 		string, 
	value 		interface{},
) *SetBuidler {
	s.getSet().(bson.M)[fmt.Sprintf("%s.$.%s",sliceField, field)] = value
	return s
}

// add to "$set" key-value in slice like singe value
// if this field exist ovetwrite value
// if you want to add field as object user ObjectBuilder
// return:
// 	{"$set": {"sliceField.$": value}}
func (s *SetBuidler) SetSliceValue(
	sliceField	string,
	value		interface{},
) *SetBuidler {
	s.getSet().(bson.M)[fmt.Sprintf("%s.$", sliceField)] = value
	return s
}


type PipelineBuilder struct {
	raw []bson.M
}

func Pipeline() *PipelineBuilder {
	return newPipelineBuilder()
}

func newPipelineBuilder() *PipelineBuilder {
	return &PipelineBuilder{
		raw: []bson.M{},
	}
}

func (p *PipelineBuilder) ToPipeline() interface{} {
	return p.raw
}

func (p *PipelineBuilder) BSON() []bson.M {
	return p.raw
}

func (p *PipelineBuilder) AddStage(
	f func(o *ObjectBuilder),
) *PipelineBuilder {
	objBuilder := Object()
	f(objBuilder)
	p.raw = append(p.raw, objBuilder.raw)
	return p
}

func (p *PipelineBuilder) MergePipelines(b *PipelineBuilder) *PipelineBuilder {
	p.raw = append(p.raw, b.raw...)
	return p
}

/* 
AddFields it's agragate stage

Example:
{ $addFields: { <newField>: <expression>, ... } }

Usage:
	Pipeline().AddFields(
		func(o *builder.ObjectBuilder) {
			o.AddField(
				"some_filed_1",
				"some_value_or_expression"
			)
		},
		func(o *builder.ObjectBuilder) {
			o.AddField(
				"some_filed_2",
				"some_value_or_expression"
			)
		}
	)
*/
func (p *PipelineBuilder) AddFields(fs ...func(o *ObjectBuilder)) *PipelineBuilder {
	o := Object()
	for _, f := range fs {
		f(o)
	}
	p.raw = append(p.raw, Object().AddField("$addFields", o).raw)
	return p
}

// Expressions
type Expression toObject

type ArrayExpression toArray

// Return array for giving arguments
// { $range: [ <start>, <end>, <non-zero step> ] }
// 
// If step zero with-out step
func Range(
	start,
	end,
	step interface{},
) Expression {
	array := Array().
		AddElem(start).
		AddElem(end)
	if step != nil {
		array.AddElem(step)
	}

	return Object().
		AddField(
			"$range", 
			array,
		)
}

type MapBody struct {
	Input	interface{}
	As		interface{}
	In		interface{}
}

func Map(
	body MapBody,
) Expression {
	buidler := Object().
			AddField(
				"input",
				body.Input,
			).AddField(
				"in",
				body.In,
			)
	if body.As != nil {
		buidler.AddField("as", body.As)
	}

	return Object().AddField("$map", buidler)
}

type FilterBody struct {
	Input 	interface{}
	As		interface{}
	Cond 	interface{}
}

func Filter(
	body FilterBody,
) Expression {
	builder := Object().
		AddField("input", body.Input).
		AddField("cond", body.Cond)

	if body.As != nil {
		builder.AddField("as", body.As)
	}

	return Object().AddField("$filter", builder)
}

type ReduceBody struct {
	Input 			interface{}
	InitialValue 	interface{}
	In				interface{}
}

func Reduce(
	body ReduceBody,
) Expression {
	b := Object().
			AddField("input", body.Input).
			AddField("initialValue", body.InitialValue).
			AddField("in", body.In)
	return Object().AddField("$reduce", b)
}

type RegexMatchBody struct {
	Input 	interface{}
	Regex 	string
	Options RegexOptions
}

func RegexMatch(body RegexMatchBody) Expression {
	return Object().AddField(
		"$regexMatch",
		Object().
			AddField(
				"input",
				body.Input,
			).
			AddField(
				"regex",
				body.Regex,
			).
			AddFieldIf(
				"options",
				option(body.Options),
				func() bool {
					return option(body.Options) != ""
				},
			),
	)
}

func Concat(
	exprs	...interface{},
) Expression {
	return Object().AddSliceField("$concat", exprs...)
}