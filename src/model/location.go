package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Project struct {
	Id   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
	Lat  string             `json:"lat,omitempty" bson:"lat,omitempty"`
	Lng  string             `json:"lng,omitempty" bson:"lng,omitempty"`
}
