package cache

import (
	"bytes"
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

func (c *Cache) decodeData(data []byte) (interface{}, error) {
	var out interface{}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Cache) Close() error {
	return nil
}

func (c *Cache) Get(id string) (interface{}, error) {
	cachedItem, err := c.client.Get(id)
	if err == nil {
		if cachedItem.Value != nil {
			return c.decodeData(cachedItem.Value)
		}
	}

	return nil, nil
}

func (c *Cache) Set(id string, value interface{}) error {
	var buff bytes.Buffer
	err := json.NewEncoder(&buff).Encode(value)
	if err != nil {
		return err
	}

	return c.client.Set(&memcache.Item{
		Key:   id,
		Value: buff.Bytes(),
	})
}

func (c *Cache) Delete(id string) error {
	return c.client.Delete(id)
}
