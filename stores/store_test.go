package store_test

import (
	"os"
	"testing"

	"github.com/amaxwellblair/api_curious/stores"
)

func TestStore_CreateUser(t *testing.T) {
	s := MustOpenStore()
	defer Close(s)

	if err := s.CreateUser("blah blah blah", "chick"); err != nil {
		t.Fatal(err)
	}

	hash, err := s.DigestToken("blah blah blah")
	if err != nil {
		t.Fatal(err)
	}

	if user, err := s.User(hash, "chick"); err != nil {
		t.Fatal(err)
	} else if user != "blah blah blah" {
		t.Fatalf("unexpected user: %s", user)
	}
}

func MustOpenStore() *store.Store {
	s := store.NewStore("../db/test_db.db")
	if err := s.Open(); err != nil {
		panic(err)
	}
	return s
}

func Close(s *store.Store) {
	s.Close()
	os.Remove("../db/test_db.db")
}
