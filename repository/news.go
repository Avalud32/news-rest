package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"news-rest/models"
)

type NewsRepository struct {
	db *sqlx.DB
}

func NewNewsRepository(db *sqlx.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

func (nr *NewsRepository) GetNewsList() ([]models.News, error) {
	var newsList []models.News
	query := `SELECT id, title, content FROM News`
	rows, err := nr.db.Queryx(query)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var news models.News
		err = rows.StructScan(&news)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		news.Categories, err = nr.getCategories(news.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get categories: %w", err)
		}
		newsList = append(newsList, news)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %w", err)
	}

	return newsList, nil
}

func (nr *NewsRepository) UpdateNews(newsId int64, news models.News) error {
	tx, err := nr.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				log.Printf("failed to rollback transaction: %v, original error: %v", errRollback, err)
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Printf("failed to commit transaction: %v", err)
			return
		}
	}()
	_, err = tx.Exec(`UPDATE News SET Title = $1, Content = $2 WHERE Id = $3`, news.Title, news.Content, newsId)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM NewsCategories WHERE NewsId = $1`, newsId)
	if err != nil {
		return fmt.Errorf("failed to delete categories: %w", err)
	}

	for _, categoryId := range news.Categories {
		_, err = tx.Exec(`INSERT INTO NewsCategories (NewsId, CategoryId) VALUES ($1, $2)`, newsId, categoryId)
		if err != nil {
			return fmt.Errorf("failed to insert categories: %w", err)
		}
	}
	return nil
}

func (nr *NewsRepository) getCategories(newsId int64) ([]int64, error) {
	var categories []int64
	query := `SELECT categoryid FROM NewsCategories WHERE newsid = $1`
	rows, err := nr.db.Queryx(query, newsId)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category int64
		err = rows.Scan(&category)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		categories = append(categories, category)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %w", err)
	}
	return categories, nil
}

func (nr *NewsRepository) updateCategories(tx *sqlx.Tx, newsId int64, categories []int64) error {
	query := `DELETE FROM NewsCategories WHERE newsid = ?`
	_, err := tx.Exec(query, newsId)
	if err != nil {
		return fmt.Errorf("failed to delete categories: %w", err)
	}

	if len(categories) == 0 {
		return nil
	}

	query = `INSERT INTO NewsCategories (newsid, categoryid) VALUES`
	var args []interface{}
	var placeholders []string
	for _, categoryId := range categories {
		placeholders = append(placeholders, "(?, ?)")
		args = append(args, newsId, categoryId)
	}
	query += " " + joinString(placeholders, ", ")
	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert categories: %w", err)
	}
	return nil
}

func joinString(elements []string, separator string) string {
	result := ""
	for i, element := range elements {
		result += element
		if i < len(elements)-1 {
			result += separator
		}
	}
	return result
}
