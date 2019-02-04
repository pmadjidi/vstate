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
	store chan *VehicleEntity
	delete chan *VehicleEntity
	Router *mux.Router
	garage *Garage
}

var app = App{wg: &sync.WaitGroup{},
	start: make(chan interface{}),
	quit:  make(chan interface{}),
	store: make(chan *VehicleEntity,100),
	delete: make(chan *VehicleEntity,100),
	port:  ":8080",
    garage: NewGarage(),
	Router: mux.NewRouter()}




