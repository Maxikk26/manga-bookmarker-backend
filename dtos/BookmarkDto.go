package dtos

type CreateBookmark struct {
	Url             string `json:"url"`
	UserId          string `json:"userId,omitempty"`
	LastChapterRead string `json:"lastChapterRead,omitempty"`
}
