package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	OK  = "OK"
	Nil = "<nil>"
)

func main() {
	var persistencePath string
	flag.StringVar(&persistencePath, "p", "./gokv-data", "Data persistence path including file name")
	flag.Parse()

	db := newDatabase()
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := db.save(persistencePath); err != nil {
					fmt.Printf("ERROR: %s\n", fmt.Sprintf("persisting database: %v", err))
				}
			}
		}
	}()

	for {
		fmt.Print(">> ")
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Printf("ERROR: %s\n", fmt.Sprintf("reading line: %v", err))
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		msg, err := db.parse(line)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			continue
		}
		fmt.Println(msg)
	}
}
