package cache

import (
	"context"
	"encoding/json"

	"github.com/bradfitz/gomemcache/memcache"
)

type Cache struct {
	client *memcache.Client
}

func NewCache(url string) (*Cache, error) {
	client := memcache.New(url)
	if err := client.Ping(); err != nil {
		return nil, err
	}

	return &Cache{client: client}, nil
}

func (c *Cache) fromJson(data []byte) (interface{}, error) {
	var out interface{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Cache) Close() error {
	return nil
}

func (c *Cache) Get(ctx context.Context, id string) (interface{}, error) {
	cachedItem, err := c.client.Get(id)
	if err == nil {
		if cachedItem.Value != nil {
			return c.fromJson(cachedItem.Value)
		}
	}

	return nil, nil
}

func (c *Cache) Set(ctx context.Context, id string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(&memcache.Item{
		Key:   id,
		Value: data,
	})
}

func (c *Cache) Delete(ctx context.Context, id string) error {
	return c.client.Delete(id)
}
