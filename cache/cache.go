package cache

import (
	"errors"

	"github.com/likejehu/urlshortener/models"
)

// ErrorNotFound is 404 err for map, when key is not found
var ErrorNotFound = errors.New("key not found")

//Storage is struct for saving links to the cache
type Storage struct {
	C map[string]*models.MineURL
}

// LinksCache is instance of cache storage
var LinksCache = NewStorage()

// CashedStore is struct for saving links to the cache
type CashedStore struct {
	C map[string]*models.MineURL
}

//NewStorage is for init new *Storage
func NewStorage() *Storage {
	var s Storage
	s.C = make(map[string]*models.MineURL, 100000)
	return &s
}

// Get fetches the specified key from the underlying store
func (s *Storage) Get(key string) (string, error) {
	if v, ok := s.C[key]; ok {
		return v.LongURL, nil
	}
	err := ErrorNotFound
	return "", err
}

// Set updates the specified key in the underlying store
func (s *Storage) Set(key string, value *models.MineURL) error {
	s.CleanUp()
	s.C[key] = value
	return nil
}

// Exists checks if the specified key is in use
// Returns false when an error occurs.
func (s *Storage) Exists(key string) (ok bool, err error) {
	if _, ok := s.C[key]; ok {
		return ok, nil
	}
	err = ErrorNotFound
	return ok, err
}

// CleanUp delete all elemets from the map in cache store
func (s *Storage) CleanUp() error {
	if len(s.C) == 99999 {
		for k := range s.C {
			delete(s.C, k)
		}
		return nil
	}
	return nil
}
