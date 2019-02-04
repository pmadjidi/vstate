package main

import (
	"database/sql"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"os"
	"sync"
)

type App struct {
	DB     *sql.DB
	port   string
	wg     *sync.WaitGroup
	quit   chan interface{}
	start  chan interface{}
	store  chan *VehicleEntity
	delete chan *VehicleEntity
	Router *mux.Router
	garage *Garage
	jwt    *jwtmiddleware.JWTMiddleware
}

var app = App{wg: &sync.WaitGroup{},
	start:  make(chan interface{}),
	quit:   make(chan interface{}),
	store:  make(chan *VehicleEntity, 100),
	delete: make(chan *VehicleEntity, 100),
	port: getPortForApp(),
	garage: NewGarage(),
	Router: mux.NewRouter(),
	jwt: jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	}),
}


func getPortForApp() string {
	res := os.Getenv("VPORT")
	if (res == "")  {
		res = ":8000"
		fmt.Printf("Setting port to default: %s\n",res)
	} else {
		res = ":" + res
		fmt.Printf("Setting port to: %s\n",res)
	}
	return res
}
