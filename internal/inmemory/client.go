package inmemory

import (
	"sync"

	"github.com/maxrasky/crema/internal/model"
)

type Client struct {
	items map[string]*model.Item
	mu    sync.RWMutex
}

func New() *Client {
	return &Client{
		items: make(map[string]*model.Item),
	}
}

func (c *Client) Set(item *model.Item) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[item.Key] = item

	return nil
}

func (c *Client) Get(key string) (*model.Item, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, ok := c.items[key]; ok {
		return item, nil
	}

	return nil, model.ErrNotFound
}

func (c *Client) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)

	return nil
}
