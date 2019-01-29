package main

import "fmt"
import "github.com/segmentio/ksuid"

func uniqueId() string {
	id := ksuid.New()
	fmt.Printf("new id: %s\n", id.String())
	return id.String()
}

