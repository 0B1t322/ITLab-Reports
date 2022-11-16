package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

type IDChecker struct {
	OnIDInvalid error
}

func NewIDChecker(
	onIDInvalid error,
) IDChecker {
	return IDChecker{
		OnIDInvalid: onIDInvalid,
	}
}

func (i IDChecker) ParseID(id string) (primitive.ObjectID, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, i.OnIDInvalid
	}

	return objId, nil
}
