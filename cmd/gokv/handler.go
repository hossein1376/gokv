package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (db *database) handleSet(cmd []string) error {
	switch len(cmd) {
	case 2:
		db.set(cmd[0], cmd[1], nil)
		return nil

	case 4:
		unit, err := parseDurationUnit(cmd[2])
		if err != nil {
			return err
		}
		duration, err := strconv.ParseUint(cmd[3], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid duration: %s", cmd[3])
		}
		expire := time.Now().Add(time.Duration(duration) * unit)

		db.set(cmd[0], cmd[1], &expire)
		return nil

	default:
		return fmt.Errorf("invalid number of arguments: %d", len(cmd))
	}
}

func (db *database) handleGet(cmd []string) error {
	switch len(cmd) {
	case 1:
		fmt.Println(db.get(cmd[0]))
		return nil

	default:
		return fmt.Errorf("invalid number of arguments: %d", len(cmd))
	}
}

func (db *database) handleSave(cmd []string) error {
	switch len(cmd) {
	case 1:
		return db.save(cmd[0])
	default:
		return fmt.Errorf("invalid number of arguments: %d", len(cmd))
	}
}

func (db *database) handleLoad(cmd []string) error {
	switch len(cmd) {
	case 1:
		return db.load(cmd[0])
	default:
		return fmt.Errorf("invalid number of arguments: %d", len(cmd))
	}
}

func parseDurationUnit(d string) (time.Duration, error) {
	switch strings.ToLower(d) {
	case "ex":
		return time.Second, nil
	case "px":
		return time.Millisecond, nil
	default:
		return 0, fmt.Errorf("invalid unit %s", d)
	}
}
