package main

import (
	"fmt"
	"net/http"

	"gitlab.tq-nest.lan/lancet/kvcache/vault"
	"gitlab.tq-nest.lan/lancet/kvcache/web"
)

func main() {
	router := new(web.Svc)
	storage := new(vault.Store)
	storage.Init(30)
	routing := []web.Route{
		{URL: "/storage/{id}", Methods: []string{"GET"}, Handler: router.GetValue},
		{URL: "/storage/{id}/{value}", Methods: []string{"PUT", "POST"}, Handler: router.SetValue},
		{URL: "/storage/{id}", Methods: []string{"DELETE"}, Handler: router.DelValue},
	}
	router.InitRouter(routing, storage.Exchange)
	http.Handle("/", router.GetRouter())
	fmt.Println("Listening")
	http.ListenAndServe(":8881", nil)
	return
}
