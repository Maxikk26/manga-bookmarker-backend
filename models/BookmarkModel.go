package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bookmark struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MangaId  primitive.ObjectID `bson:"mangaId"`
	UserId   primitive.ObjectID `bson:"userId"`
	Chapter  string             `bson:"chapter"`
	LastRead primitive.DateTime `bson:"lastRead"`
	Status   int                `bson:"status"`
}
