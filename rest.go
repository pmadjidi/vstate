package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
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
	stateName := args["state"]
	st, ok := vstate.ValidState(stateName)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("<h1>Not Found</>"))
	} else {
		v, ok := app.garage.getVehicleById(id)
		if (ok) {
			req := &Request{vstate.SetState, vstate.Admins, st, make(chan *Response)}
			v.Port <- req
			var res *Response
			select {
			case res = <-req.resp:
			case <-time.After(5 * time.Second):
				w.WriteHeader(http.StatusRequestTimeout)
				w.Write([]byte("<h1>Timeout....</>"))
				return
			}

			if (res.err == nil) {
				w.WriteHeader(http.StatusOK)
				b, err := json.Marshal(v)
				if (err == nil) {
					w.Write(b)
				} else {
					w.WriteHeader(http.StatusExpectationFailed)
					w.Write([]byte(res.err.Error()))
					return
				}
			}
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<h1>Not Found</>"))
	}
}

func (a *App) claimVehicle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Req: %s %s\n", r.URL.Host, r.URL.Path)
	args := mux.Vars(r)
	id := args["id"]
	v, ok := app.garage.getVehicleById(id)
	if (ok) {
		req := &Request{vstate.Claim, vstate.Admins, vstate.Nothing, make(chan *Response)}
		v.Port <- req

		var res *Response
		select {
		case res = <-req.resp:
		case <-time.After(5 * time.Second):
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("<h1>Timeout....</>"))
			return
		}

		if (res.err == nil) {
			w.WriteHeader(http.StatusOK)
			b, err := json.Marshal(v)
			if (err == nil) {
				w.Write(b)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("<h>Ok...  but no content....</>"))
				return
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("<h1>" + res.err.Error() + "</>"))
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("<h1>Not Found</>"))
}

func init() {
	fmt.Print("Configuring the rest api...\n")
	app.initializeRoutes()
}
