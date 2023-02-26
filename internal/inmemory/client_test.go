package inmemory

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/maxrasky/crema/internal/model"
)

func TestClient_Operations(t *testing.T) {
	client := New()

	key := "InMemory_Unique_Key_#$%"

	tests := []struct {
		name     string
		callback func() (*model.Item, error)
		expected *model.Item
		err      error
	}{
		{
			name: "get non-stored key",
			callback: func() (*model.Item, error) {
				return client.Get(key)
			},
			err: model.ErrNotFound,
		},
		{
			name: "set value",
			callback: func() (*model.Item, error) {
				newItem := &model.Item{
					Key:   key,
					Value: []byte(`dancing in the dark`),
				}

				err := client.Set(newItem)
				return nil, err
			},
		},
		{
			name: "get value",
			callback: func() (*model.Item, error) {
				return client.Get(key)
			},
			expected: &model.Item{
				Key:   key,
				Value: []byte(`dancing in the dark`),
			},
		},
		{
			name: "delete key",
			callback: func() (*model.Item, error) {
				err := client.Delete(key)
				return nil, err
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			item, err := tt.callback()
			if tt.err != nil {
				assert.ErrorIs(t, err, tt.err)
				return
			}

			assert.NoError(t, err)
			if tt.expected != nil {
				assert.NotNil(t, item)
				assert.Equal(t, *tt.expected, *item)
			}
		})
	}
}
