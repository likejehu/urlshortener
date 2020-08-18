package db

import (
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v7"
	"github.com/likejehu/urlshortener/models"
)

// Error404 is 404 err for store, when key is not found
var Error404 = errors.New("key not found")

// RedisStore is an implementation of the Storer Interface backed by Redis
type RedisStore struct {
	client *redis.Client
}

// LinksStore is instance of  local redis storage with default port 6379
var LinksStore = NewRedisStore("localhost:6379")

// Get fetches the specified key from the underlying store
func (r *RedisStore) Get(key string) (string, error) {
	s, err := r.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			err = Error404
			return s, err
		}
		return s, err
	}
	b := []byte(s)
	m := &models.MineURL{}
	err = json.Unmarshal(b, m)
	if err != nil {
		return s, err
	}
	return m.LongURL, err
}

// Set updates the specified key in the underlying store
func (r *RedisStore) Set(key string, value *models.MineURL) error {
	var v interface{} = value
	return r.client.Set(key, v, 0).Err()
}

// Exists checks if the specified key is in use
// Returns false when an error occurs.
func (r *RedisStore) Exists(key string) (bool, error) {
	_, err := r.client.Exists(key).Result()
	if err != nil {
		return false, err
	}

	return true, nil
}

// NewRedisStore returns a new Store Instance
func NewRedisStore(address string) *RedisStore {
	var store = &RedisStore{}
	store.client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return store
}
