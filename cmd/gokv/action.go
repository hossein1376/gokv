package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type database struct {
	DB map[string]value
	mu sync.RWMutex
}

type value struct {
	Value     string
	CreatedAt time.Time
	Expire    *time.Time
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

func (db *database) get(key string) string {
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
		return ""
	}

	if v.Expire == nil {
		return v.Value
	}
	if v.Expire.Before(time.Now()) {
		func() {
			db.mu.Lock()
			defer db.mu.Unlock()
			delete(db.DB, key)
		}()
		return ""
	}
	return v.Value
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

func (db *database) parse(line string) error {
	if line == "" {
		return nil
	}
	switch cmd := strings.Fields(line); strings.ToLower(cmd[0]) {
	case "exit":
		os.Exit(0)
		return nil

	case "help":
		printUsage()
		return nil

	case "set":
		return db.handleSet(cmd[1:])

	case "get":
		return db.handleGet(cmd[1:])

	case "save":
		return db.handleSave(cmd[1:])

	case "load":
		return db.handleLoad(cmd[1:])

	default:
		return fmt.Errorf("unknown command: %s", cmd[0])
	}
}

func printUsage() {
	msg := `Usage: gokv <command> [<args>]
commands:
	set <key> <value> [<unit> <duration>]
	get <key>
	save <file>
	load <file>
	help
	exit
duration unit:
	ex = seconds
	px = milliseconds
Note: keywords are case insensitive.`
	fmt.Println(msg)
}
