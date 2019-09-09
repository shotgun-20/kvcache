package vault

import (
	"time"
)

// cleaner - уничтожаем устаревшие ключи
func (store *Store) cleaner() {
	reply := make(chan Message)
	for {
		msg := Message{Action: "POP", Reply: reply}
		store.exchange <- msg
		got := <-reply
		if got.Error == true {
			time.Sleep(1 * time.Second)
		}
	}
}

// control - управление хранилищем
func (store *Store) control() {
	tick := time.Tick(time.Second) // Каждую секунду ставим в очередь разделитель "вёдер"
	for {
		select {
		case <-tick:
			store.addNode("", "", false)
			for {
				err := store.popNode()
				if err != nil {
					break
				}
			}
		case request := <-store.exchange:
			reply := Message{}
			switch request.Action {
			/*
				case "POP": // внутренний метод, часть механизма устаревания. Недоступен клиентам.
					err := store.popNode()
					if err != nil {
						reply.Error = true
					}
			*/
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
			request.Reply <- reply
		}
	}
}
