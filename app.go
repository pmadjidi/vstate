package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"sync"
)



type App struct {
	DB    *sql.DB
	port  string
	wg    *sync.WaitGroup
	quit  chan interface{}
	start chan interface{}
	store chan *Vehicle
	delete chan *Vehicle
	Router *mux.Router
	garage *Garage
}

var app = App{wg: &sync.WaitGroup{},
	start: make(chan interface{}),
	quit:  make(chan interface{}),
	store: make(chan *Vehicle,100),
	delete: make(chan *Vehicle,100),
	port:  ":8080",
    garage: NewGarage(),
	Router: mux.NewRouter()}




