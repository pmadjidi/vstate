package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"time"
	"vehicles/vstate"
)

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/admin/newv", a.createVehicle).Methods("GET")
	a.Router.HandleFunc("/admin/listv", a.listVehicle).Methods("GET")
	a.Router.HandleFunc("/admin/claim/{id:[0-9A-Za-z]+}", a.claimDisClaimVehicle).Methods("GET")
	a.Router.HandleFunc("/admin/disclaim/{id:[0-9]+}", a.claimDisClaimVehicle).Methods("GET")
	a.Router.HandleFunc("/admin/setstate/{id:[0-9A-Za-z]+}/{state:[a-zA-z]+}", a.setState).Methods("Get")
	//a.Router.HandleFunc("/admin/getstate/{id:[0-9A-Za-z]+}", a.getState).Methods("GET")
	//a.Router.HandleFunc("/admin/delete/{id:[0-9A-Za-z]+}", a.deleteVehicle).Methods("DELETE")

	// a.Router.HandleFunc("/hunter/hunt/{id:[0-9]+}", a.hunt).Methods("GET")
	a.Router.HandleFunc("/admin/claim/{id:[0-9]+}", a.claimDisClaimVehicle).Methods("GET")
	a.Router.HandleFunc("/admin/disclaim/{id:[0-9]+}", a.claimDisClaimVehicle).Methods("GET")

	a.Router.HandleFunc("/user/claim/{id:[0-9A-Za-z]+}", a.claimDisClaimVehicle).Methods("GET")
	a.Router.HandleFunc("/user/disclaim/{id:[0-9A-Za-z]+}", a.claimDisClaimVehicle).Methods("GET")

}

func (a *App) createVehicle(w http.ResponseWriter, r *http.Request) {
	v := NewVehicle()
	app.garage.Set(v.Uid, v)
	b, err := json.Marshal(v)
	if (err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out("Error...", 2))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (a *App) listVehicle(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(app.garage.getMap())
	if (err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out("Error...", 2))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

var validPath = regexp.MustCompile("^/(admin|user|hunter)/(a-zA-Z)+/.+$")

func (a *App) setState(w http.ResponseWriter, r *http.Request) {
	role := validPath.FindStringSubmatch(r.URL.Path)
	fmt.Printf("Req: %s %s %s\n", r.URL.Host, r.URL.Path,role[2])
	args := mux.Vars(r)
	id := args["id"]
	stateName := args["state"]
	st, ok := vstate.ValidState(stateName)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write(out("NotFound...", 2))
	} else {
		v, ok := app.garage.getVehicleById(id)
		if (ok) {
			req := &Request{vstate.SetState, vstate.Admins, st, make(chan *Response)}
			v.Port <- req
			select {
			case res := <-req.resp:
				if (res.err == nil) {
					w.WriteHeader(http.StatusOK)
					b, err := json.Marshal(v)
					if (err == nil) {
						w.Write(b)
					} else {
						w.WriteHeader(http.StatusExpectationFailed)
						w.Write(out(res.err.Error(), 2))
						return
					}
				} else {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(wraphtml(res.err.Error(), 1)))
					return
				}
			case <-time.After(5 * time.Second):
				w.WriteHeader(http.StatusRequestTimeout)
				w.Write(out("Timeout...", 2))
				return
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write(out("NotFound", 2))
		}
	}
}


func (a *App) claimDisClaimVehicle(w http.ResponseWriter, r *http.Request) {
	parse := validPath.FindStringSubmatch(r.URL.Path)
	roleName := parse[0]
	eventName := parse[1]
	var role vstate.URoles
	var event vstate.Event

	switch roleName {
	case "admin":
		role = vstate.Admins
	case "hunter":
		role = vstate.Hunters
	case "user":
		role = vstate.EndUser
	default:
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(wraphtml("Forbidden...", 1)))
		return
	}

	switch eventName {
	case "claim":
		event = vstate.Claim
	case "disclaim":
		event = vstate.DisClaim
	default:
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(wraphtml("Forbidden...", 1)))
		return
	}

	args := mux.Vars(r)
	id := args["id"]
	v, ok := app.garage.getVehicleById(id)
	if (ok) {
		req := &Request{event, role, vstate.Nothing, make(chan *Response)}
		v.Port <- req
		select {
		case res := <-req.resp:
			if (res.err == nil) {
				w.WriteHeader(http.StatusOK)
				b, err := json.Marshal(v)
				if (err == nil) {
					w.Write(b)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write(out("OK, but not content...", 2))
					return
				}
			} else {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(wraphtml(res.err.Error(), 1)))
				return
			}
		case <-time.After(5 * time.Second):
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write(out("TimeOut...", 2))
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write(out("NotFound...", 2))
	}
}

func init() {
	fmt.Print("Configuring the rest api...\n")
	app.initializeRoutes()
}
