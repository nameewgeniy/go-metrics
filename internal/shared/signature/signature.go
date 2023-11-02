package signature

import (
	"crypto/sha256"
	"fmt"
	"io"
)

type Signature interface {
	Hash(data []byte) (string, error)
	Valid(hash string, data []byte) (bool, error)
}

var Sign Signature

type ChekerSignature struct {
	key string
}

func Singleton(key string) error {

	if key == "" {
		return fmt.Errorf("ChekerSignature: Singleton: key cannot be empty")
	}

	Sign = &ChekerSignature{
		key: key,
	}

	return nil
}

func (s ChekerSignature) Hash(data []byte) (string, error) {
	message := string(data) + s.key
	h := sha256.New()

	if _, err := io.WriteString(h, message); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (s ChekerSignature) Valid(hash string, data []byte) (bool, error) {
	calculatedHash, err := s.Hash(data)
	if err != nil {
		return false, fmt.Errorf("valid: An error occurred while validating the signature")
	}

	return calculatedHash == hash, nil
}
