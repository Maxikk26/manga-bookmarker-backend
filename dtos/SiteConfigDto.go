package dtos

type CreateSiteConfig struct {
	Name            string `json:"name"`
	BaseUrl         string `json:"baseUrl"`
	TitleSelector   string `json:"titleSelector"`
	ChapterSelector string `json:"chapterSelector"`
	CoverSelector   string `json:"coverSelector"`
	UploadSelector  string `json:"uploadSelector"`
	GenreSelector   string `json:"genreSelector,omitempty"`
}