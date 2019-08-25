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
	reply := make(chan vault.Message)
	msg := vault.Message{Action: "GET", Reply: reply, Key: id}
	svc.exchange <- msg
	resp := <-reply
	if resp.Error != true {
		fmt.Fprint(w, resp.Value)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "FAILURE")
}

// SetValue - установить/обновить значение ключа
func (svc *Svc) SetValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	v := vars["value"]
	reply := make(chan vault.Message)
	msg := vault.Message{Action: "SET", Reply: reply, Key: id, Value: v}
	svc.exchange <- msg
	resp := <-reply
	if resp.Error != true {
		fmt.Fprint(w, "OK")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "FAILURE")
}

// DelValue - удалить ключ из хранилища
func (svc *Svc) DelValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	reply := make(chan vault.Message)
	msg := vault.Message{Action: "DEL", Reply: reply, Key: id}
	svc.exchange <- msg
	resp := <-reply
	if resp.Error != true {
		fmt.Fprint(w, "OK")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "FAILURE")
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
	svc.exchange = exch
	return
}

// GetRouter - возвращаем ссылку на роутер для веб-сервера
func (svc *Svc) GetRouter() *mux.Router {
	if svc.router != nil {
		return svc.router
	}
	return nil
}
