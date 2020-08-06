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
	value, err := svc.store.GetValue(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "FAILURE")
		return
	}
	fmt.Fprint(w, value)
}

// SetValue - установить/обновить значение ключа
func (svc *Svc) SetValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	v := vars["value"]
	err := svc.store.SetValue(id, v)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "FAILURE")
		return
	}
	fmt.Fprint(w, "OK")
}

// DelValue - удалить ключ из хранилища
func (svc *Svc) DelValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := svc.store.DelValue(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "FAILURE")
		return
	}
	fmt.Fprint(w, "OK")
}

// InitRouter - инициализировать маршрутизацию запросов
func (svc *Svc) InitRouter(routes []Route, store *vault.Store) {
	if svc.router != nil {
		return
	}
	svc.router = mux.NewRouter()
	for _, i := range routes {
		for _, j := range i.Methods {
			svc.router.HandleFunc(i.URL, i.Handler).Methods(j)
		}
	}
	svc.store = store
	return
}

// GetRouter - возвращаем ссылку на роутер для веб-сервера
func (svc *Svc) GetRouter() *mux.Router {
	if svc.router != nil {
		return svc.router
	}
	return nil
}
