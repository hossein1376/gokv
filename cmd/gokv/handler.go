package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func (db *database) parse(line string) (string, error) {
	switch cmd := strings.Fields(line); strings.ToLower(cmd[0]) {
	case "exit":
		os.Exit(0)
		return "", nil

	case "help":
		return handleUsage(), nil

	case "set":
		return db.handleSet(cmd[1:])

	case "get":
		return db.handleGet(cmd[1:])

	case "save":
		err := db.handleSave(cmd[1:])
		if err != nil {
			return "", err
		}
		return OK, nil

	case "load":
		err := db.handleLoad(cmd[1:])
		if err != nil {
			return "", err
		}
		return OK, nil

	default:
		return "", fmt.Errorf("unknown command: %s", cmd[0])
	}
}

func (db *database) handleSet(cmd []string) (string, error) {
	switch len(cmd) {
	case 2:
		db.set(cmd[0], cmd[1], nil)
		return OK, nil

	case 4:
		unit, err := parseDurationUnit(cmd[2])
		if err != nil {
			return "", err
		}
		duration, err := strconv.ParseUint(cmd[3], 10, 64)
		if err != nil {
			return "", fmt.Errorf("invalid duration: %s", cmd[3])
		}
		expire := time.Now().Add(time.Duration(duration) * unit)

		db.set(cmd[0], cmd[1], &expire)
		return OK, nil

	default:
		return "", fmt.Errorf("invalid number of arguments: %d", len(cmd))
	}
}

func (db *database) handleGet(cmd []string) (string, error) {
	switch len(cmd) {
	case 1:
		v, found := db.get(cmd[0])
		if !found {
			return Nil, nil
		}
		return v, nil

	default:
		return "", fmt.Errorf("invalid number of arguments: %d", len(cmd))
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

func handleUsage() string {
	return `Usage: gokv <command> [<args>]
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
