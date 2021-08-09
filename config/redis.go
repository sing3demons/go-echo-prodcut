package config

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/sing3demons/go-echo-product/models"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

type RedisCache interface {
	Set(key string, value interface{})
	GetProduct(key string) []models.Products
	GetPage(key string) interface{}
}

func NewRedisCache(host string, db int, exp time.Duration) RedisCache {
	return &redisCache{host: host, db: db, expires: exp}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value interface{}) {
	client := cache.getClient()

	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	client.Set(key, json, cache.expires*time.Second)
}
func (cache *redisCache) GetProduct(key string) []models.Products {
	client := cache.getClient()

	val, err := client.Get(key).Result()
	if err != nil {
		return nil
	}

	product := []models.Products{}
	err = json.Unmarshal([]byte(val), &product)
	if err != nil {
		panic(err)
	}

	return product
}

func (cache *redisCache) GetPage(key string) interface{} {

	client := cache.getClient()

	val, err := client.Get(key).Result()
	if err != nil {
		return nil
	}

	return val

}
