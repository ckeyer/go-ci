package types

import (
	"gopkg.in/mgo.v2/bson"
)

type Project struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Name    string        `json:"name" bson:"name"`
	Type    string        `json:"type" bson:"type"`
	Members []Member      `json:"members" bson:"members"`
}

type Member struct {
	UserId bson.ObjectId `json:"user_id" bson:"user_id"`
	Role   string        `json:"role" bson:"role"`
}
