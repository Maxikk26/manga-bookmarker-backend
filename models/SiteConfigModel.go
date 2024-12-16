package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SiteConfig struct {
	Id              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UpdatedAt       primitive.DateTime `bson:"updatedAt"`
	Name            string             `bson:"name"`
	BaseUrl         string             `bson:"baseUrl"`
	TitleSelector   string             `bson:"titleSelector"`
	ChapterSelector string             `bson:"chapterSelector"`
	CoverSelector   string             `bson:"coverSelector"`
	UploadSelector  string             `bson:"uploadSelector"`
	GenreSelector   string             `bson:"genreSelector"`
}
