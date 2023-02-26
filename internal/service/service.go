package service

import (
	"context"

	"github.com/maxrasky/crema/internal/service/proto"

	"github.com/maxrasky/crema/internal/model"
)

type Storager interface {
	Get(key string) (*model.Item, error)
	Set(item *model.Item) error
	Delete(key string) error
}

type Service struct {
	store Storager
}

func New(store Storager) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Get(_ context.Context, key *proto.Key) (*proto.Item, error) {
	item, err := s.store.Get(key.Key)
	if err != nil {
		return nil, err
	}

	return &proto.Item{
		Key:   item.Key,
		Value: item.Value,
	}, nil
}

func (s *Service) Set(_ context.Context, item *proto.Item) (*proto.Null, error) {
	if err := s.store.Set(&model.Item{
		Key:   item.Key,
		Value: item.Value,
	}); err != nil {
		return nil, err
	}

	return &proto.Null{}, nil
}

func (s *Service) Delete(_ context.Context, key *proto.Key) (*proto.Null, error) {
	if err := s.store.Delete(key.Key); err != nil {
		return nil, err
	}

	return &proto.Null{}, nil
}
