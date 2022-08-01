package hash

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const saltLength = 8

// ErrIncorrectPassword is returned when the provided password is incorrect.
var ErrIncorrectPassword = errors.New("password is not correct")

// Argon2Hasher uses Argon2 to hash passwords with random salt.
type Argon2Hasher struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

const (
	defaultArgon2Time    uint32 = 1
	defaultArgon2Memory  uint32 = 64 * 1024
	defaultArgon2Threads uint8  = 4
	defaultArgon2KeyLen  uint32 = 32
)

// Option is a function that can be used to customize the Argon2Hasher.
type Option func(*Argon2Hasher)

// Time sets the custom time parameter of the Argon2 algorithm.
func Time(time uint32) Option {
	return func(ah *Argon2Hasher) {
		ah.time = time
	}
}

// Memory sets the custom memory parameter of the Argon2 algorithm.
func Memory(memory uint32) Option {
	return func(ah *Argon2Hasher) {
		ah.memory = memory
	}
}

// Threads sets the custom threads parameter of the Argon2 algorithm.
func Threads(threads uint8) Option {
	return func(ah *Argon2Hasher) {
		ah.threads = threads
	}
}

// KeyLen sets the custom key length parameter of the Argon2 algorithm.
func KeyLen(keyLen uint32) Option {
	return func(ah *Argon2Hasher) {
		ah.keyLen = keyLen
	}
}

// NewArgon2Hasher creates a new Argon2Hasher.
func NewArgon2Hasher(opts ...Option) Argon2Hasher {
	hasher := Argon2Hasher{
		time:    defaultArgon2Time,
		memory:  defaultArgon2Memory,
		threads: defaultArgon2Threads,
		keyLen:  defaultArgon2KeyLen,
	}

	for _, opt := range opts {
		opt(&hasher)
	}

	return hasher
}

func (hasher Argon2Hasher) hashPassword(password, salt []byte) []byte {
	hashedPassword := argon2.IDKey(password, salt, hasher.time, hasher.memory, hasher.threads, hasher.keyLen)

	return append(salt, hashedPassword...)
}

// Hash returns the argon2 hash of the password.
func (hasher Argon2Hasher) Hash(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hashedPassword := hasher.hashPassword([]byte(password), salt)

	return string(hashedPassword), nil
}

// Check checks if the provided password is correct or not.
func (hasher Argon2Hasher) Check(password, hashedPassword string) error {
	hash := []byte(hashedPassword)
	salt := make([]byte, saltLength)
	copy(salt, hash)

	ok := bytes.Equal(hasher.hashPassword([]byte(password), salt), hash)
	if !ok {
		return ErrIncorrectPassword
	}

	return nil
}
