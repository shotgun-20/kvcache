package vault

import (
	"time"
)

// control - управление хранилищем
func (store *Store) control() {
	tick := time.NewTimer(store.ttl)
	for {
		select {
		case <-tick.C: // запускаем устаревание элементов
			now := time.Now()
			step := store.ttl
			for {
				if store.head != nil {
					diff := store.head.expire.Sub(now)
					if diff < 1*time.Second {
						// узел протух, убираем его
						//fmt.Printf("pop! @%s\n", now)
						err := store.popNode()
						if err != nil || store.head == nil {
							//fmt.Print("nothing more to pop\n")
							step = store.ttl
							break
						}
						continue
					} else {
						step = store.head.expire.Sub(now)
						break
					}
				} else {
					break
				}
			}
			tick = time.NewTimer(step)
		case request, ok := <-store.exchange:
			if !ok {
				break
			}
			reply := Message{}
			switch request.Action {
			case "SET":
				//fmt.Printf("SET: %s\n", request.Key)
				err := store.setNode(request.Key, request.Value)
				reply.Key = request.Key
				if err != nil {
					reply.Error = true
				}
			case "GET":
				//fmt.Printf("GET: %s\n", request.Key)
				v, err := store.getNode(request.Key)
				if err != nil {
					reply.Error = true
				} else {
					reply.Key = request.Key
					reply.Value = v
				}
			case "DEL":
				//fmt.Printf("DEL: %s\n", request.Key)
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
