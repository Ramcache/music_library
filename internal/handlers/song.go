package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"music_library/internal/config"
	"music_library/internal/models"
)

type Handler struct {
	DB     *gorm.DB
	Config *config.Config
	Log    *logrus.Logger
}

type Pagination struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=10"`
}

// GetSongs
// @Summary Получение всех песен
// @Description Получение списка всех песен с возможностью фильтрации и пагинации
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Фильтр по группе"
// @Param song query string false "Фильтр по названию песни"
// @Param release_date query string false "Фильтр по дате выпуска"
// @Param link query string false "Фильтр по ссылке"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество элементов на странице"
// @Success 200 {object} map[string]interface{}
// @Router /songs [get]
func (h *Handler) GetSongs(c *gin.Context) {
	var pagination Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := h.DB.Model(&models.Song{})

	filters := c.Request.URL.Query()
	for key, values := range filters {
		if key != "page" && key != "limit" {
			query = query.Where(key+" = ?", values[0])
		}
	}

	var songs []models.Song
	var total int64
	query.Count(&total)
	query.Offset(offset).Limit(pagination.Limit).Find(&songs)

	c.JSON(http.StatusOK, gin.H{
		"data":  songs,
		"total": total,
		"page":  pagination.Page,
		"limit": pagination.Limit,
	})
}

// GetSongText
// @Summary Получение текста песни
// @Description Получение текста песни по её ID с пагинацией
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество элементов на странице"
// @Success 200 {object} map[string]interface{}
// @Router /songs/{id}/text [get]
func (h *Handler) GetSongText(c *gin.Context) {
	id := c.Param("id")
	var song models.Song
	if result := h.DB.First(&song, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
		return
	}

	verses := strings.Split(song.Text, "\n\n")
	var pagination Pagination
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offset := (pagination.Page - 1) * pagination.Limit
	end := offset + pagination.Limit
	if end > len(verses) {
		end = len(verses)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  verses[offset:end],
		"total": len(verses),
		"page":  pagination.Page,
		"limit": pagination.Limit,
	})
}

// DeleteSong
// @Summary Удаление песни
// @Description Удаление песни по её ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 204
// @Router /songs/{id} [delete]
func (h *Handler) DeleteSong(c *gin.Context) {
	id := c.Param("id")
	if result := h.DB.Delete(&models.Song{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateSong
// @Summary Обновление информации о песне
// @Description Обновление данных о песне по её ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param input body models.Song true "Обновлённые данные песни"
// @Success 200 {object} models.Song
// @Router /songs/{id} [put]
func (h *Handler) UpdateSong(c *gin.Context) {
	id := c.Param("id")
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Model(&models.Song{}).Where("id = ?", id).Updates(song); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, song)
}

// @Summary Добавление новой песни
// @Description Добавление новой песни с данными из внешнего API
// @Tags songs
// @Accept json
// @Produce json
// @Param input body models.Song true "Данные новой песни"
// @Success 201 {object} models.Song
// @Router /songs [post]
func (h *Handler) AddSong(c *gin.Context) {
	var req struct {
		Group string `json:"group" binding:"required"`
		Song  string `json:"song" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiURL := h.Config.APIBaseURL + "/info"
	resp, err := http.Get(apiURL + "?group=" + req.Group + "&song=" + req.Song)
	if err != nil {
		h.Log.WithError(err).Error("Ошибка при запросе к внешнему API")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка внешнего API"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "Внешний API вернул ошибку"})
		return
	}

	var detail struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		h.Log.WithError(err).Error("Ошибка декодирования ответа")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка"})
		return
	}

	song := models.Song{
		Group:       req.Group,
		Song:        req.Song,
		ReleaseDate: detail.ReleaseDate,
		Text:        detail.Text,
		Link:        detail.Link,
	}

	if result := h.DB.Create(&song); result.Error != nil {
		h.Log.WithError(result.Error).Error("Ошибка создания записи в БД")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}

	c.JSON(http.StatusCreated, song)
}
