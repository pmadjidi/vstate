package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"strings"
	"time"
	"vehicles/vstate"
     "github.com/swaggo/http-swagger"
	_"vehicles/docs"
)

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/admin/newv", createVehicle).Methods("GET")
	a.Router.HandleFunc("/admin/listv", listVehicle).Methods("GET")
	a.Router.HandleFunc("/admin/claim/{id:[0-9A-Za-z]+}", claimDisClaimVehicle).Methods("GET")
	a.Router.HandleFunc("/admin/disclaim/{id:[0-9A-Za-z]+}", claimDisClaimVehicle).Methods("GET")
	a.Router.HandleFunc("/admin/setstate/{id:[0-9A-Za-z]+}/{state:[a-zA-z]+}", setState).Methods("Get")
	a.Router.HandleFunc("/admin/getstate/{id:[0-9A-Za-z]+}", getState).Methods("GET")
	a.Router.HandleFunc("/admin/delete/{id:[0-9A-Za-z]+}", deleteVehicle).Methods("GET")

	a.Router.HandleFunc("/hunter/hunt/{id:[0-9A-Za-z]+}", huntVehicle).Methods("GET")
	a.Router.HandleFunc("/hunter/claim/{id:[0-9A-Za-z]+}", claimDisClaimVehicle).Methods("GET")
	a.Router.HandleFunc("/hunter/disclaim/{id:[0-9A-Za-z]+}", claimDisClaimVehicle).Methods("GET")

	a.Router.HandleFunc("/user/claim/{id:[0-9A-Za-z]+}", claimDisClaimVehicle).Methods("GET")
	a.Router.HandleFunc("/user/disclaim/{id:[0-9A-Za-z]+}", claimDisClaimVehicle).Methods("GET")
	a.Router.HandleFunc("/swagger/*", httpSwagger.WrapHandler)


}

func  createVehicle(w http.ResponseWriter, r *http.Request) {
	v := NewVehicle()
	app.garage.Set(v.Uid, v)
	b, err := v.serilize()
	if (err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out("Error...", 2))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func deleteVehicle(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	id := args["id"]
	v, ok := app.garage.getVehicleById(id)
	if (!ok) {
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(http.StatusNotFound)
		w.Write(out("NotFound", 2))
		return
	} else {
		req := &Request{vstate.Delete, vstate.Admins, vstate.Nothing, make(chan *Response)}
		v.Port <- req
		_ = <-req.resp
		app.delete <- v
		w.WriteHeader(http.StatusOK)
		w.Write(out("Deleted", 2))
		return
	}
}

func huntVehicle(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	id := args["id"]
	v, ok := app.garage.getVehicleById(id)
	if (!ok) {
		w.Header().Set("Content-Type", "application/text")
		w.WriteHeader(http.StatusNotFound)
		w.Write(out("NotFound", 2))
		return
	} else {
		req := &Request{vstate.Hunter, vstate.Hunters, vstate.Nothing, make(chan *Response)}
		v.Port <- req

		select {
		case res := <-req.resp:
			if (res.err == nil) {
				w.WriteHeader(http.StatusOK)
				b, err := v.serilize()
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

		w.WriteHeader(http.StatusOK)
		w.Write(out("Deleted", 2))
		return
	}
}

func listVehicle(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(app.garage.getMap())
	if (err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(out("Error...", 2))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

var validPath = regexp.MustCompile("^/(admin|user|hunter)/([a-zA-Z0-9]+)/.+$")

func setState(w http.ResponseWriter, r *http.Request) {
	role := validPath.FindStringSubmatch(r.URL.Path)
	fmt.Printf("Req: %s %s %s\n", r.URL.Host, r.URL.Path, role)
	args := mux.Vars(r)
	id := args["id"]
	stateName := strings.Title(args["state"])

	st, ok := vstate.ValidState(stateName)
	if !ok {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write(out("NotFound...", 2))
	} else {
		fmt.Printf("setGetState: Got state %s\n", st.String())
		v, ok := app.garage.getVehicleById(id)
		if (ok) {
			req := &Request{vstate.SetState, vstate.Admins, st, make(chan *Response)}
			v.Port <- req
			select {
			case res := <-req.resp:
				if (res.err == nil) {
					w.WriteHeader(http.StatusOK)
					b, err := v.serilize()
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

func getState(w http.ResponseWriter, r *http.Request) {
	role := validPath.FindStringSubmatch(r.URL.Path)
	fmt.Printf("Req: %s %s %s\n", r.URL.Host, r.URL.Path, role)
	args := mux.Vars(r)
	id := args["id"]

	v, ok := app.garage.getVehicleById(id)
	if (ok) {
		w.WriteHeader(http.StatusOK)
		b, err := v.serilize()
		if (err == nil) {
			w.Write(b)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(out("OK, but not content...", 2))
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write(out("NotFound", 2))
	}
}

func claimDisClaimVehicle(w http.ResponseWriter, r *http.Request) {
	parse := validPath.FindStringSubmatch(r.URL.Path)
	roleName := parse[1]
	eventName := parse[2]
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
		fmt.Printf("id:  %s Event: %s Role: %s \n",id,event,role.String())
		req := &Request{event, role, vstate.Nothing, make(chan *Response)}
		v.Port <- req
		select {
		case res := <-req.resp:
			if (res.err == nil) {
				w.WriteHeader(http.StatusOK)
				b, err := v.serilize()
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
		fmt.Printf("id: %s not found\n",id)
		w.WriteHeader(http.StatusNotFound)
		w.Write(out("NotFound...", 2))
	}
}

func init() {
	fmt.Print("Configuring the rest api...\n")
	app.initializeRoutes()
}
