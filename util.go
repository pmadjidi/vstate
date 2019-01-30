package main

import (
	"fmt"
	"log"
)
import "github.com/segmentio/ksuid"

func uniqueId() string {
	id := ksuid.New()
	fmt.Printf("new id: %s\n", id.String())
	return id.String()
}

func exit(err error) {
	log.Fatal(err)
	close(app.quit)
}

