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

// @Summary Get all songs
// @Description Get all songs with filtering and pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Group filter"
// @Param song query string false "Song filter"
// @Param release_date query string false "Release date filter"
// @Param link query string false "Link filter"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
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

// @Summary Get song lyrics
// @Description Get paginated song lyrics by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /songs/{id}/text [get]
func (h *Handler) GetSongText(c *gin.Context) {
	id := c.Param("id")
	var song models.Song
	if result := h.DB.First(&song, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
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

// @Summary Delete a song
// @Description Delete a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
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

// @Summary Update a song
// @Description Update a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param input body models.Song true "Song data"
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

// @Summary Add a new song
// @Description Add a new song with data from external API
// @Tags songs
// @Accept json
// @Produce json
// @Param input body models.Song true "Song data"
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

	// Call external API
	apiURL := h.Config.APIBaseURL + "/info"
	resp, err := http.Get(apiURL + "?group=" + req.Group + "&song=" + req.Song)
	if err != nil {
		h.Log.WithError(err).Error("Failed to call external API")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "External API error"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "External API returned error"})
		return
	}

	var detail struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		h.Log.WithError(err).Error("Failed to decode response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
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
		h.Log.WithError(result.Error).Error("Failed to create song")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusCreated, song)
}
