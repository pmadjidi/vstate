package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)


func main() {
	fmt.Print("Starting service....\n")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for _ = range c {
			fmt.Print("\nStopping the service....\n")
			app.DB.Close()
			close(app.quit)
		}
	}()

	go func () {
		srv := &http.Server{
			Handler:      app.Router,
			Addr:         "127.0.0.1:8000",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		close(app.start)
		log.Fatal(srv.ListenAndServe())
	}()

	<-app.quit

}
