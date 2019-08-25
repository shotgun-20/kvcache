package vault

import (
	"errors"
	"fmt"
	"time"
)

// Init - инициализация хранилища
func (store *store) Init(ttl uint64) error {
	if store.ttl != 0 {
		return errors.New("Already up")
	}
	if ttl < 1 {
		return errors.New("Wrong TTL")
	}

	store.ttl = ttl
	store.Exchange = make(chan Message)
	store.flat = make(map[string]*node)

	var i uint64
	for i = 0; i < ttl; i++ {
		store.addNode("", "", false)
	}
	go store.control() // запускаем управление хранилищем
	go store.cleaner() // запускаем устаревание записей
	return nil
}

func (store *store) cleaner() {
	reply := make(chan Message)
	msg := Message{Action: "POP", Reply: reply}
	for {
		store.Exchange <- msg
		got := <-reply
		if got.Error == true {
			fmt.Print("Cleaner sleeps.\n")
			time.Sleep(1 * time.Second)
		}
	}
}

// control - управление хранилищем
func (store *store) control() {
	tick := time.Tick(time.Second) // Каждую секунду ставим в очередь разделитель "вёдер"
	for {
		select {
		case <-tick:
			fmt.Print("Got time tick\n")
			store.addNode("", "", false)
		case request := <-store.Exchange:
			fmt.Print("Got request", request, "\n")
			reply := new(Message)
			switch request.Action {
			case "POP": // внутренний метод, часть механизма устаревания. Недоступен клиентам.
				err := store.popNode()
				if err != nil {
					reply.Error = true
				}
			case "SET":
				err := store.setNode(request.Key, request.Value)
				reply.Key = request.Key
				if err != nil {
					reply.Error = true
				}
			case "GET":
				v, err := store.getNode(request.Key)
				if err != nil {
					reply.Error = true
				} else {
					reply.Key = request.Key
					reply.Value = v
				}
			case "DEL":
				err := store.delNode(request.Key)
				if err != nil {
					reply.Key = request.Key
					reply.Error = true
				}
			}
		}
	}
}
