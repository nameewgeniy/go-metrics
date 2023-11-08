package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

type Signature interface {
	Hash(data []byte) string
	Valid(hash string, data []byte) bool
}

var Sign Signature

type ChekerSignature struct {
	key string
}

func Singleton(key string) {
	if Sign != nil {
		return
	}

	Sign = &ChekerSignature{
		key: key,
	}
}

func (s ChekerSignature) Hash(data []byte) string {
	h := hmac.New(sha256.New, []byte(s.key))
	h.Write(data)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s ChekerSignature) Valid(hash string, data []byte) bool {
	calculatedHash := s.Hash(data)

	return calculatedHash == hash
}
