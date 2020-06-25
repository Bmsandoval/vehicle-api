package cache_service

import (
	"bytes"
	"encoding/gob"
)

func (s ServiceImplementation) Set(key string, value interface{}) error {
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	// Encode (send) the value.
	err := enc.Encode(value)
	if err != nil {
		return err
	}

	network.Bytes()

	s.Cache.Set(key, network.Bytes(), s.Timeout)

	return nil
}
