package jwt

import "errors"

type Key struct {
	bytes  []byte
	method string
	id     string
}

func NewKey(bytes []byte, method string, id string) (Key, error) {
	if len(bytes) == 0 {
		return Key{}, errors.New("key bytes cannot be empty")
	}
	if method == "" {
		return Key{}, errors.New("key method cannot be empty")
	}
	if id == "" {
		return Key{}, errors.New("key id cannot be empty")
	}
	return Key{bytes: bytes, method: method, id: id}, nil
}

func (k Key) Bytes() []byte  { return k.bytes }
func (k Key) Method() string { return k.method }
func (k Key) ID() string     { return k.id }

func (k Key) IsRSA() bool {
	switch k.method {
	case "RS256", "RS384", "RS512":
		return true
	}
	return false
}

func (k Key) IsECDSA() bool {
	switch k.method {
	case "ES256", "ES384", "ES512":
		return true
	}
	return false
}
