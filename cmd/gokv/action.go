package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"
	"time"
)

type database struct {
	DB map[string]value
	mu sync.RWMutex
}

type value struct {
	Value  string
	Expire *time.Time
}

func newDatabase() *database {
	return &database{
		DB: make(map[string]value),
	}
}

func (db *database) set(k, v string, expire *time.Time) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.DB[k] = value{
		Value:  v,
		Expire: expire,
	}
}

func (db *database) get(key string) (string, bool) {
	var (
		v  value
		ok bool
	)
	func() {
		db.mu.RLock()
		defer db.mu.RUnlock()
		v, ok = db.DB[key]
	}()
	if !ok {
		return "", false
	}

	if v.Expire == nil {
		return v.Value, true
	}
	if v.Expire.Before(time.Now()) {
		func() {
			db.mu.Lock()
			defer db.mu.Unlock()
			delete(db.DB, key)
		}()
		return "", false
	}
	return v.Value, true
}

func (db *database) save(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	db.mu.RLock()
	defer db.mu.RUnlock()
	enc := gob.NewEncoder(f)
	if err = enc.Encode(db.DB); err != nil {
		return fmt.Errorf("encode database: %w", err)
	}
	return nil
}

func (db *database) load(path string) error {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	db.mu.RLock()
	defer db.mu.RUnlock()
	dec := gob.NewDecoder(f)
	if err = dec.Decode(&db.DB); err != nil {
		return fmt.Errorf("decode data into database: %w", err)
	}
	return nil
}
