package port

import "github.com/bitly/domain"

type MappingService interface {
	Shorten(longUrl string) (*domain.ShortenResult, error)
	Click(mappingId uint, ipAddress, userAgent string) (string, error)
}
