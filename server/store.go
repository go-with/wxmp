package server

import (
	"errors"
)

var (
	ErrTokenNotFound          = errors.New("token is not found in store")
	ErrEncodingAESKeyNotFound = errors.New("encoding aes key is not found in store")
)

type Store interface {
	GetToken(appID string) (token string, err error)
	GetEncodingAESKey(appID string) (encodingAESKey string, err error)
}

type MemoryStore struct {
	tokens map[string]string
	keys   map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tokens: make(map[string]string),
		keys:   make(map[string]string),
	}
}

func (sto *MemoryStore) GetToken(appID string) (token string, err error) {
	token, ok := sto.tokens[appID]
	if !ok {
		err = ErrTokenNotFound
	}
	return
}

func (sto *MemoryStore) GetEncodingAESKey(appID string) (encodingAESKey string, err error) {
	encodingAESKey, ok := sto.keys[appID]
	if !ok {
		err = ErrEncodingAESKeyNotFound
	}
	return
}

func (sto *MemoryStore) SetToken(appID, token string) {
	sto.tokens[appID] = token
}

func (sto *MemoryStore) SetEncodingAESKey(appID, encodingAESKey string) {
	sto.keys[appID] = encodingAESKey
}
