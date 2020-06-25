package cache_service

import (
	"errors"
)

func (s ServiceImplementation) Get(key string) ([]byte, error) {
	cachedResult, ok := s.Cache.Get(key)
	if ! ok {
		// record not in cache
		return nil, nil
	}
	byteArray, ok := cachedResult.([]byte)
	if ! ok {
		// invalid record in cache
		return nil, errors.New("non-byte record found in cache")
	}

	return byteArray, nil
}
