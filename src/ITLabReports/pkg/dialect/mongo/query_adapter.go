package mongo

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ArgNotFound 	= errors.New("Arg not found")
	UknownOperation = errors.New("Uknown operation")
)

type opearation int
const (
	InsertOne	opearation = iota
	InsertMany
	DeleteOne
	DeleteMany
	UpdateByID
	UpdateOne
	UpdateMany
	ReplaceOne
	Drop

	Aggragate	
	CountDocuments
	EstimatedDocumentCount
	Distinct
	Find
	FindOne
	FindOneAndDelete
	FindOneAndReplace
	FindOneAndUpdate
	Watch
)


type Query struct {
	sc				mongo.SessionContext
	Collection		string						`json:"collection" 	bson:"collection"`
	Operation		opearation					`json:"operaion" 	bson:"operation"`
	Args			map[string]interface{}		`json:"args" 		bson:"args"`
}

func (q Query) String() string {
	data, _ := json.Marshal(q)
	buf := bytes.NewBuffer(data)
	return buf.String()
}

func (q *Query) query(coll *mongo.Collection) (*mongo.Cursor, error) {
	switch q.Operation {
	case Aggragate:
		return q.agragate(coll)
	case Find:
		return q.find(coll)
	}
	return nil, UknownOperation
}

func (q *Query) agragate(coll *mongo.Collection) (*mongo.Cursor, error) {
	pipeline, err := q.mustGetArg("pipeline")
	if err != nil {
		return nil, err
	}

	return coll.Aggregate(
		q.sc,
		pipeline,
	)
}

func (q *Query) find(coll *mongo.Collection) (*mongo.Cursor, error) {
	filter, err := q.mustGetArg("filter")
	if err != nil {
		return nil, err
	}

	return coll.Find(
		q.sc,
		filter,
	)
}

func (q *Query) mustGetArg(Arg string) (interface{}, error) {
	arg, find := q.Args[Arg]
	if !find {
		return nil, errors.Wrap(ArgNotFound, Arg)
	}

	return arg, nil
}

func (q *Query) Query(
	dbName	string,
) (*mongo.Cursor, error) {
	coll := q.sc.Client().
				Database(dbName).
				Collection(q.Collection)

	return q.query(coll)
}

func (q *Query) Exec(
	dbName	string,
) (interface{}, error) {
	coll := q.sc.Client().
				Database(dbName).
				Collection(q.Collection)
	return q.exec(coll)
}

func (q *Query) exec(coll *mongo.Collection) (interface{}, error) {
	switch q.Operation {
	}
	return nil, UknownOperation
}

func QueryFromString(
	sc mongo.SessionContext, 
	query string,
) *Query {
	buf := bytes.NewBufferString(query)
	q := Query{}
	json.Unmarshal(buf.Bytes(), &q)
	q.sc = sc

	return &q
}