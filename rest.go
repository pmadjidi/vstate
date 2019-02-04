package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"net/http"
	"regexp"
	"strings"
	"time"
	_ "vehicles/docs"
	"vehicles/vstate"
)



func (a *App) jwtHandler(f func(w http.ResponseWriter,r *http.Request)) http.HandlerFunc  {
	return func(w http.ResponseWriter,r *http.Request) {
		// user := context.Get(r, "user")
		//fmt.Fprintf(w, "This is an authenticated request")
		//fmt.Fprintf(w, "Claim content:\n")
		/*
		for k, v := range user.(*jwt.Token).Claims {
			fmt.Fprintf(w, "%s :\t%#v\n", k, v)
		}
		*/
		f(w,r)
	}
}

func (a *App) initializeRoutes() {

	a.Router.HandleFunc("/admin/newv", createVehicle).Methods("GET")
	a.Router.HandleFunc("/admin/listv", a.jwtHandler(listVehicle)).Methods("GET")
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
	//a.Router.HandleFunc("/swagger/*", httpSwagger.WrapHandler)
	a.Router.PathPrefix("/doc/").Handler(httpSwagger.WrapHandler)

}



func  createVehicle(w http.ResponseWriter, r *http.Request) {
	v := NewVehicle()
	app.garage.Set(v.Uid, v)
	b, err := v.serilize()
	if (err != nil) {
		w.WriteHeader(http.StatusInternalServerError)
		Respond(w,http.StatusInternalServerError,"Error....")
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func deleteVehicle(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	id := args["id"]
	v, ok := app.garage.getVehicleById(id)
	if (!ok) {
		Respond(w,http.StatusNotFound,"NotFound....")
	} else {
		req := &Request{vstate.Delete, vstate.Admins, vstate.Nothing, make(chan *Response)}
		v.Port <- req
		_ = <-req.resp
		v.delete()
		Respond(w,http.StatusOK,"Deleted....")
		return
	}
}

func huntVehicle(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	id := args["id"]
	v, ok := app.garage.getVehicleById(id)
	if (!ok) {
		Respond(w,http.StatusNotFound, "NotFound....")
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
					Respond(w,http.StatusExpectationFailed, "No data....")
				}
			} else {
				Respond(w,http.StatusForbidden, res.err.Error())
			}
		case <-time.After(5 * time.Second):
			Respond(w,http.StatusRequestTimeout, "Timeout....")
		}
	}
}

func listVehicle(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(app.garage.getMap())
	if (err != nil) {
		Respond(w,http.StatusInternalServerError,"Error....")
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
		Respond(w,http.StatusNotAcceptable,"State not acceptable....")
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
						Respond(w,http.StatusExpectationFailed,res.err.Error())
					}
				} else {
					Respond(w,http.StatusForbidden,res.err.Error())
				}
			case <-time.After(5 * time.Second):
				Respond(w,http.StatusRequestTimeout,"Timeout....")
			}
		} else {
			Respond(w,http.StatusNotFound,"NotFound....")
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
			Respond(w,http.StatusExpectationFailed,"No data....")
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		Respond(w,http.StatusNotFound,"Not Found")
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
		Respond(w,http.StatusForbidden,"Forbidden....")
	}

	switch eventName {
	case "claim":
		event = vstate.Claim
	case "disclaim":
		event = vstate.DisClaim
	default:
		Respond(w,http.StatusForbidden,"Forbidden....")
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
					Respond(w,http.StatusExpectationFailed,"No data....")
				}
			} else {
				Respond(w,http.StatusForbidden,res.err.Error())
			}
		case <-time.After(5 * time.Second):
			Respond(w,http.StatusRequestTimeout,"Timeout....")
		}
	} else {
		Respond(w,http.StatusNotFound,"NotFound....")
	}
}

func init() {
	fmt.Print("Configuring the rest api...\n")
	app.initializeRoutes()
}
