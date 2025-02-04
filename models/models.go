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

func (n *News) Validate() error {
	if n.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if n.Content == "" {
		return fmt.Errorf("content cannot be empty")
	}
	return nil
}
