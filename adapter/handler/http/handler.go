package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/bitly/domain"
	"github.com/bitly/usecase"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	MappingService usecase.MappingService
}

func NewHandler(mappingService usecase.MappingService) *Handler {
	return &Handler{MappingService: mappingService}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/shorten", h.ShortenHandler)
	r.GET("/:code", h.RedirectHandler)
}

type ShortenRequest struct {
	LongUrl string `json:"long_url" binding:"required,url"`
}

type ShortenResponse struct {
	ShortUrl  string `json:"short_url"`
	ShortCode string `json:"short_code"`
}

func (h *Handler) ShortenHandler(c *gin.Context) {
	var payload ShortenRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "long_url is required and must be a valid URL"})
		return
	}

	result, err := h.MappingService.Shorten(payload.LongUrl)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ShortenResponse{
		ShortUrl:  baseUrl(c) + "/" + result.ShortCode,
		ShortCode: result.ShortCode,
	})
}

func (h *Handler) RedirectHandler(c *gin.Context) {
	longUrl, err := h.MappingService.Click(c.Param("code"), c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, longUrl)
}

// baseUrl builds "http(s)://host" from the incoming request, e.g. "http://localhost:8080".
func baseUrl(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}
