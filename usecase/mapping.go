package usecase

import (
	"time"

	"github.com/bitly/domain"
	"github.com/bitly/port"
)

type MappingService struct {
	MappingRepository port.MappingRepository
	LinkTTL           time.Duration
}

func NewMappingService(mappingRepository port.MappingRepository, linkTTL time.Duration) *MappingService {
	return &MappingService{MappingRepository: mappingRepository, LinkTTL: linkTTL}
}

func (m *MappingService) Shorten(longUrl string) (*domain.ShortenResult, error) {
	urlMapping, err := m.MappingRepository.FindByLongURL(longUrl)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(m.LinkTTL)

	if urlMapping != nil {
		if urlMapping.ExpiresAt != nil && urlMapping.ExpiresAt.Compare(time.Now()) <= 0 {
			return m.MappingRepository.Renew(urlMapping.ID, &expiresAt)
		}

		return &domain.ShortenResult{ID: urlMapping.ID, ShortCode: urlMapping.ShortCode}, nil
	}

	shortenResult, err := m.MappingRepository.Shorten(longUrl, &expiresAt)

	if err != nil {
		return nil, err
	}

	return shortenResult, nil
}

func (m *MappingService) Click(shortCode, ipAddress, userAgent string) (string, error) {
	urlMapping, err := m.MappingRepository.FindByShortCode(shortCode)
	if err != nil {
		return "", err
	}

	if urlMapping == nil {
		return "", domain.ErrNotFound
	}

	if urlMapping.ExpiresAt != nil && urlMapping.ExpiresAt.Compare(time.Now()) <= 0 {
		return "", domain.ErrNotFound
	}

	err = m.MappingRepository.Log(urlMapping.ID, ipAddress, userAgent)
	if err != nil {
		return "", err
	}

	return urlMapping.LongUrl, nil
}
