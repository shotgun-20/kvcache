package web

import (
	"fmt"
	"net/http"

	"gitlab.tq-nest.lan/lancet/kvcache/vault"

	"github.com/gorilla/mux"
)

// GetValue - получить значение ключа из хранилища
func (svc *Svc) GetValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	resp := fmt.Sprintf("Get ID: %s\n", id)
	fmt.Fprint(w, resp)
}

// SetValue - установить/обновить значение ключа
func (svc *Svc) SetValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	resp := fmt.Sprintf("Set ID: %s\n", id)
	fmt.Fprint(w, resp)
}

// DelValue - удалить ключ из хранилища
func (svc *Svc) DelValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	resp := fmt.Sprintf("Delete ID: %s\n", id)
	fmt.Fprint(w, resp)
}

// InitRouter - инициализировать маршрутизацию запросов
func (svc *Svc) InitRouter(routes []Route, exch chan vault.Message) {
	if svc.router != nil {
		return
	}
	svc.router = mux.NewRouter()
	for _, i := range routes {
		for _, j := range i.Methods {
			svc.router.HandleFunc(i.URL, i.Handler).Methods(j)
		}
	}
	svc.Exchange = exch
	return
}

// GetRouter - возвращаем ссылку на роутер для веб-сервера
func (svc *Svc) GetRouter() *mux.Router {
	if svc.router != nil {
		return svc.router
	}
	return nil
}
