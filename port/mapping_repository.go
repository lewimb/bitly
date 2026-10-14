package port

import (
	"time"

	"github.com/bitly/domain"
)

type MappingRepository interface {
	Shorten(longUrl string, expiresAt *time.Time) (*domain.ShortenResult, error)
	Renew(mappingId uint, expiresAt *time.Time) (*domain.ShortenResult, error)
	FindByLongURL(longUrl string) (*domain.URLMapping, error)
	FindByShortCode(shortCode string) (*domain.URLMapping, error)
	Log(mappingId uint, ipAddress, userAgent string) error
}
