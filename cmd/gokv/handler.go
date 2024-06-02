package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

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
