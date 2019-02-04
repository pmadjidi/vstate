package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func wraphtml(s string,h int) string {
	s = strings.TrimSuffix(s, "\n")
	res :=  "<h" + strconv.Itoa(h) + ">" + s + "</>"
	fmt.Print(res + "\n")
	return  res
}

func out(s string,size int) []byte {
	return []byte(wraphtml(s,size))
}

type Jmessage map[string]interface{}

func Message(status int, message string) (Jmessage) {
	return map[string]interface{} {"status" : http.StatusText(status), "message" : message}
}

func Respond(w http.ResponseWriter, code int,msg string)  {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Message(code,msg))
}


