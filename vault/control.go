package vault

import (
	"time"
)

// control - управление хранилищем
func (store *Store) control() {
	tick := time.Tick(time.Second) // Каждую секунду
	// ставим в очередь разделитель "вёдер" и убираем устаревшее
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
		case request, ok := <-store.exchange:
			if !ok {
				break
			}
			reply := Message{}
			switch request.Action {
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
