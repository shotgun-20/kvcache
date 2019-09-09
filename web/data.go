package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.tq-nest.lan/lancet/kvcache/vault"
)

// Route - правило маршрутизации запросов
type Route struct {
	URL     string
	Methods []string
	Handler func(http.ResponseWriter, *http.Request)
}

// Svc - обмен данными с хранилищем
type Svc struct {
	store  *vault.Store
	router *mux.Router
}
