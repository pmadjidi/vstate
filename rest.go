package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"vehicles/vstate"
)

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/admin/newv", a.createVehicle).Methods("GET")
	a.Router.HandleFunc("/admin/listv", a.listVehicle).Methods("GET")
	a.Router.HandleFunc("/admin/claim/{id:[0-9A-Za-z]+}", a.claimVehicle).Methods("GET")
	/*
	a.Router.HandleFunc("/admin/disclaim/{id:[0-9]+}", a.disclaimVehicle).Methods("GET")
	a.Router.HandleFunc("/admin/setstate/{id:[0-9]+}/{state:[a-zA-z]+}", a.setState).Methods("PUT")
	a.Router.HandleFunc("/admin/getstate/{id:[0-9]+}", a.getState).Methods("GET")
	a.Router.HandleFunc("/admin/delete/{id:[0-9]+}", a.deleteVehicle).Methods("DELETE")

	a.Router.HandleFunc("/hunter/hunt/{id:[0-9]+}", a.hunt).Methods("GET")
	a.Router.HandleFunc("/admin/claim/{id:[0-9]+}", a.claimVehicle).Methods("POST")
	a.Router.HandleFunc("/admin/disclaim/{id:[0-9]+}", a.disclaimVehicle).Methods("GET")


	a.Router.HandleFunc("/user/claim/{id:[0-9]+}", a.claimVehicle).Methods("GET")
	a.Router.HandleFunc("/user/disclaim/{id:[0-9]+}", a.disclaimVehicle).Methods("PUT")
	*/
}

func (a *App) createVehicle(w http.ResponseWriter, r *http.Request) {
	v := NewVehicle()
	app.garage.Set(v.Uid, v)
	b, err := json.Marshal(v)
	if (err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<h1>Error 500</>"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (a *App) listVehicle(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(app.garage.getMap())
	if (err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<h1>Error 500</>"))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (a *App) setState(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Req: %s %s\n", r.URL.Host, r.URL.Path)
	args := mux.Vars(r)
	id := args["id"]
	state := args["state"]
	switch(state) {
	case: ""
	}
	v, ok := app.garage.getVehicleById(id)
	if (ok) {
		v.Set()
		w.WriteHeader(http.StatusOK)
		b, err := json.Marshal(v)
		if (err == nil) {
			w.Write(b)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("<h1>Not Found</>"))
		}

	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>OK</>"))
}

func (a *App) claimVehicle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Req: %s %s\n", r.URL.Host, r.URL.Path)
	args := mux.Vars(r)
	id := args["id"]
	v, ok := app.garage.getVehicleById(id)
	if (ok) {
		v.Port <- Request{vstate.Claim, vstate.Admins}
		w.WriteHeader(http.StatusOK)
		b, err := json.Marshal(v)
		if (err == nil) {
			w.Write(b)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("<h1>Not Found</>"))
		}

	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<h1>OK</>"))
	}
}

func init() {
	fmt.Print("Configuring the rest api...\n")
	app.initializeRoutes()
}
