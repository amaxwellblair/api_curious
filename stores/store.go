package store

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/boltdb/bolt"
)

// Store holds the key/value store
type Store struct {
	Path string
	DB   *bolt.DB
}

// NewStore returns a new instance of a Store
func NewStore(path string) *Store {
	return &Store{
		Path: path,
	}
}

// Open opens a database file or creates one
func (s *Store) Open() error {
	var err error
	if s.DB, err = bolt.Open(s.Path, 0600, &bolt.Options{Timeout: 1 * time.Second}); err != nil {
		return err
	}

	// Initialize users bucket
	if s.DB.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte("users")); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// Close closes the database that was previously opened
func (s *Store) Close() error {
	return s.DB.Close()
}

// CreateUser creates a user for a given token
func (s *Store) CreateUser(token, salt string) error {
	// Hash and salt the user's token
	hash, err := s.DigestToken(token)
	if err != nil {
		return err
	}
	saltedHash, err := s.DigestToken(hash + salt)
	if err != nil {
		return err
	}

	// Create or Overwrite user
	if err := s.DB.Update(func(tx *bolt.Tx) error {
		tx.Bucket([]byte("users")).Put([]byte(saltedHash), []byte(token))
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// User retrieves a user from the key/value store
func (s *Store) User(hash, salt string) (string, error) {
	var token string
	saltedHash, err := s.DigestToken(hash + salt)
	if err != nil {
		return "", err
	}
	s.DB.View(func(tx *bolt.Tx) error {
		token = string(tx.Bucket([]byte("users")).Get([]byte(saltedHash)))
		return nil
	})
	return token, nil
}

// DigestToken hashes and salts a given token
func (s *Store) DigestToken(value string) (string, error) {
	hashAlgo := sha256.New()
	if _, err := hashAlgo.Write([]byte(value)); err != nil {
		return "", err
	}
	hash := hashAlgo.Sum(nil)

	// Encode to hexadecmial for brevity and compatibility
	hexHash := hex.EncodeToString(hash)

	return string(hexHash), nil
}
