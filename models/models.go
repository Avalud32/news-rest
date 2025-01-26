package models

type News struct {
	Id         int64   `json:"Id" db:"id"`
	Title      string  `json:"Title" db:"title"`
	Content    string  `json:"Content" db:"content"`
	Categories []int64 `json:"Categories,omitempty"`
}

type NewsCategories struct {
	NewsId     int64 `db:"newsid"`
	CategoryId int64 `db:"categoryid"`
}

type NewsResponse struct {
	Success bool   `json:"Success"`
	News    []News `json:"News"`
	Error   string `json:"Error,omitempty"`
}
