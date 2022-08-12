// Package hash implements the password hashing algorithms.
package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// ErrIncorrectPassword is returned when the provided password is incorrect.
var ErrIncorrectPassword = errors.New("password is not correct")

// Argon2Hasher uses Argon2 to hash passwords with random salt.
type Argon2Hasher struct {
	format  string
	version int
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

const (
	defaultArgon2Time    uint32 = 1
	defaultArgon2Memory  uint32 = 64 * 1024
	defaultArgon2Threads uint8  = 4
	defaultArgon2KeyLen  uint32 = 32
	defaultArgon2SaltLen uint32 = 32
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

// SaltLen sets the custom salt length parameter of the Argon2 algorithm.
func SaltLen(saltLen uint32) Option {
	return func(ah *Argon2Hasher) {
		ah.saltLen = saltLen
	}
}

// NewArgon2Hasher creates a new Argon2Hasher.
func NewArgon2Hasher(opts ...Option) Argon2Hasher {
	hasher := Argon2Hasher{
		format:  "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		version: argon2.Version,
		time:    defaultArgon2Time,
		memory:  defaultArgon2Memory,
		threads: defaultArgon2Threads,
		keyLen:  defaultArgon2KeyLen,
		saltLen: defaultArgon2SaltLen,
	}

	for _, opt := range opts {
		opt(&hasher)
	}

	return hasher
}

// Hash returns the argon2 hash of the password.
func (ah Argon2Hasher) Hash(plain string) (string, error) {
	salt := make([]byte, ah.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(plain), salt, ah.time, ah.memory, ah.threads, ah.keyLen)

	return fmt.Sprintf(
		ah.format,
		ah.version,
		ah.memory,
		ah.time,
		ah.threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

// Check checks if the provided password is correct or not.
func (ah Argon2Hasher) Check(plain, hash string) error {
	hashParts := strings.Split(hash, "$")

	_, err := fmt.Sscanf(hashParts[3], "m=%d,t=%d,p=%d", &ah.memory, &ah.time, &ah.threads)
	if err != nil {
		return fmt.Errorf("failed to parse hash: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(hashParts[4])
	if err != nil {
		return fmt.Errorf("failed to decode salt: %w", err)
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(hashParts[5])
	if err != nil {
		return fmt.Errorf("failed to decode hash: %w", err)
	}

	hashToCompare := argon2.IDKey([]byte(plain), salt, ah.time, ah.memory, ah.threads, uint32(len(decodedHash)))

	if subtle.ConstantTimeCompare(hashToCompare, decodedHash) == 1 {
		return nil
	}

	return ErrIncorrectPassword
}
