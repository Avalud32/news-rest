package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"news-rest/models"
	"news-rest/repository"
	"strconv"
)

type NewsHandler struct {
	newsRepo *repository.NewsRepository
}

func NewNewsHandler(repo *repository.NewsRepository) *NewsHandler {
	return &NewsHandler{
		newsRepo: repo,
	}
}

func (h *NewsHandler) GetNewsList(c *fiber.Ctx) error {
	newsList, err := h.newsRepo.GetNewsList()
	if err != nil {
		log.Printf("failed to get news list: %v", err)
		return c.JSON(models.NewsResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to get news list: %v", err),
		})
	}

	return c.JSON(models.NewsResponse{
		Success: true,
		News:    newsList,
	})
}

func (h *NewsHandler) UpdateNews(c *fiber.Ctx) error {
	idParam := c.Params("Id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewsResponse{
			Success: false,
			Error:   fmt.Sprintf("invalid id: %s", idParam),
		})
	}
	var news models.News
	if err := c.BodyParser(&news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewsResponse{
			Success: false,
			Error:   fmt.Sprintf("invalid json: %v", err),
		})
	}
	if err := h.newsRepo.UpdateNews(id, news); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewsResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to update news: %v", err),
		})
	}

	return c.JSON(models.NewsResponse{
		Success: true,
	})
}
