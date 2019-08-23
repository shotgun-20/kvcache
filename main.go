package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URL     string
	Methods []string
	Handler func(http.ResponseWriter, *http.Request)
}

func getValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	resp := fmt.Sprintf("Get ID: %s\n", id)
	fmt.Fprint(w, resp)
}

func setValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	resp := fmt.Sprintf("Set ID: %s\n", id)
	fmt.Fprint(w, resp)
}

func delValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	resp := fmt.Sprintf("Delete ID: %s\n", id)
	fmt.Fprint(w, resp)
}

func initRouter(router *mux.Router, routes []Route) {
	for _, i := range routes {
		for _, j := range i.Methods {
			router.HandleFunc(i.URL, i.Handler).Methods(j)
		}
	}
	return
}

func main() {
	router := mux.NewRouter()
	routing := []Route{
		{"/storage/{id}", []string{"GET"}, getValue},
		{"/storage/{id}", []string{"PUT", "POST"}, setValue},
		{"/storage/{id}", []string{"DELETE"}, delValue},
	}
	initRouter(router, routing)
	http.Handle("/", router)
	fmt.Println("Listening")
	http.ListenAndServe(":8881", nil)
	return
}
