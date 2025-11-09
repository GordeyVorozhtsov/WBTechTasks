package shortener

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"net/url"
	"sync"
	"time"

	"github.com/wb-go/wbf/redis"
)

const (
	alphabet   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // выборка для короткой сслыки
	codeLength = 6                                                                // длина новой(короткой) ссылки
	ttl        = 5 * time.Minute                                                  // TTL для кэша
)

type Urls struct {
	OldURL string
	NewURL string
}

type Visit struct {
	Timestamp time.Time
	UserAgent string
}

type Service struct {
	mu        sync.RWMutex
	urls      map[string]*Urls
	analytics map[string][]Visit
	cache     *redis.Client
}

func NewService() *Service {
	// подключаем Redis
	client := redis.New("localhost:6379", "", 0) // для теста — redis:alpine на localhost:6379

	return &Service{
		urls:      make(map[string]*Urls),
		analytics: make(map[string][]Visit),
		cache:     client,
	}
}

// генерируем короткую ссылку
func generateCode() (string, error) {
	b := make([]byte, codeLength)
	for i := 0; i < codeLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		b[i] = alphabet[n.Int64()]
	}
	return string(b), nil
}

// Создание новой короткой ссылки
func (s *Service) Shorten(ctx context.Context, longURL string) (string, error) {
	// валидация URL
	parsed, err := url.ParseRequestURI(longURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", errors.New("invalid url")
	}

	code, err := generateCode()
	if err != nil {
		return "", err
	}

	newURL := "/s/" + code

	s.mu.Lock()
	s.urls[code] = &Urls{OldURL: longURL, NewURL: newURL}
	s.mu.Unlock()

	// пишем в Redis с TTL
	_ = s.cache.Set(ctx, code, longURL)
	_ = s.cache.Client.Expire(ctx, code, ttl).Err()

	return newURL, nil
}

// Получение оригинальной ссылки
func (s *Service) Resolve(ctx context.Context, code string) (string, error) {
	// сначала пробуем из Redis
	val, err := s.cache.Get(ctx, code)
	if err == nil && val != "" {
		return val, nil
	}

	// fallback из памяти
	s.mu.RLock()
	defer s.mu.RUnlock()
	if u, ok := s.urls[code]; ok {
		return u.OldURL, nil
	}
	return "", errors.New("not found")
}

// Запись аналитики
func (s *Service) RecordVisit(code, userAgent string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.analytics[code] = append(s.analytics[code], Visit{
		Timestamp: time.Now(),
		UserAgent: userAgent,
	})
}

// Получение аналитики
func (s *Service) GetAnalytics(code string) ([]Visit, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if visits, ok := s.analytics[code]; ok {
		return visits, nil
	}
	return nil, errors.New("no analytics")
}
