package model

import "errors"

type Item struct {
	Key        string
	Value      []byte
	Flags      uint32
	Expiration uint32
}

var (
	ErrNotFound = errors.New("item not found")
)
