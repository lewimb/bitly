package repository

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/bitly/domain"
	"gorm.io/gorm"
)

const base62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const (
	shortCodeLength     = 6
	maxGenerateAttempts = 5
)

func GenerateShortCode(n int) (string, error) {
	b := make([]byte, n)
	max := big.NewInt(int64(len(base62)))
	for i := range b {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = base62[idx.Int64()]
	}
	return string(b), nil
}

type MappingRepository struct {
	DB  *gorm.DB
	Ctx context.Context
}

func NewMappingRepository(db *gorm.DB, ctx context.Context) *MappingRepository {
	return &MappingRepository{DB: db, Ctx: ctx}
}

func (m *MappingRepository) Shorten(longUrl string, expiresAt *time.Time) (*domain.ShortenResult, error) {
	shortCode, err := m.uniqueShortCode()
	if err != nil {
		return nil, err
	}

	result := gorm.WithResult()
	urlMapping := domain.URLMapping{LongUrl: longUrl, ShortCode: shortCode, ExpiresAt: expiresAt}

	err = gorm.G[domain.URLMapping](m.DB, result).Create(m.Ctx, &urlMapping)

	if err != nil {
		return nil, err
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("failed to create url mapping: no rows affected")
	}

	return &domain.ShortenResult{ID: urlMapping.ID, ShortCode: urlMapping.ShortCode}, nil
}

// Renew gives an expired mapping a new short code and a new expiration.
func (m *MappingRepository) Renew(mappingId uint, expiresAt *time.Time) (*domain.ShortenResult, error) {
	shortCode, err := m.uniqueShortCode()
	if err != nil {
		return nil, err
	}

	result := m.DB.WithContext(m.Ctx).
		Model(&domain.URLMapping{}).
		Where("id = ?", mappingId).
		Updates(map[string]any{
			"short_code": shortCode,
			"expires_at": expiresAt,
		})

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("failed to renew url mapping: no rows affected")
	}

	return &domain.ShortenResult{ID: mappingId, ShortCode: shortCode}, nil
}

// uniqueShortCode generates a short code that is not used by any existing mapping.
func (m *MappingRepository) uniqueShortCode() (string, error) {
	for range maxGenerateAttempts {
		shortCode, err := GenerateShortCode(shortCodeLength)
		if err != nil {
			return "", err
		}

		existing, err := m.FindByShortCode(shortCode)
		if err != nil {
			return "", err
		}

		if existing == nil {
			return shortCode, nil
		}
	}

	return "", fmt.Errorf("failed to generate a unique short code after %d attempts", maxGenerateAttempts)
}

func (m *MappingRepository) FindByLongURL(longUrl string) (*domain.URLMapping, error) {

	urlMap, err := gorm.G[domain.URLMapping](m.DB).Where("long_url = ?", longUrl).First(m.Ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &urlMap, nil
}

func (m *MappingRepository) FindByShortCode(shortCode string) (*domain.URLMapping, error) {
	urlMap, err := gorm.G[domain.URLMapping](m.DB).Where("short_code = ?", shortCode).First(m.Ctx)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &urlMap, nil
}

func (m *MappingRepository) Log(mappingId uint, ipAddress, userAgent string) error {

	result := gorm.WithResult()
	urlLog := domain.URLClick{MappingID: mappingId, IPAddress: ipAddress, UserAgent: userAgent}

	err := gorm.G[domain.URLClick](m.DB, result).Create(m.Ctx, &urlLog)

	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to log clicked url: no rows affected")
	}

	return nil
}
