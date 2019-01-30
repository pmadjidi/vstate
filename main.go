package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)


func main() {
	fmt.Print("Starting service....\n")
	close(app.start)
	srv := &http.Server{
		Handler:      app.Router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
	<-app.quit
}
