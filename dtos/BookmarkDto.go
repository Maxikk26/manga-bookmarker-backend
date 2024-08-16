package dtos

import (
	"time"
)

type CreateBookmark struct {
	Url     string `json:"url"`
	UserId  string `json:"userId,omitempty"`
	Chapter string `json:"chapter,omitempty"`
}

type Bookmark struct {
	Id              string    `json:"id"`
	Chapter         string    `json:"chapter"`
	LastChapterRead time.Time `jon:"LastChapterRead"`
}
