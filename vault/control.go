package vault

import (
	"errors"
	"time"
)

// Init - инициализация хранилища
func (store *Store) Init(ttl uint64) error {
	if store.ttl != 0 {
		return errors.New("Already up")
	}
	if ttl < 1 {
		return errors.New("Wrong TTL")
	}

	store.ttl = ttl
	store.Exchange = make(chan Message)
	store.flat = make(map[string]*node)

	// Создаём достаточное количество узлов-филлеров, чтобы механизм
	// устаревания не начал уничтожать записи преждевременно.
	var i uint64
	for i = 0; i < ttl; i++ {
		store.addNode("", "", false)
	}

	go store.control() // запускаем управление хранилищем
	go store.cleaner() // запускаем устаревание записей
	return nil
}

// cleaner - уничтожаем устаревшие ключи
func (store *Store) cleaner() {
	reply := make(chan Message)
	for {
		msg := Message{Action: "POP", Reply: reply}
		store.Exchange <- msg
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
		case request := <-store.Exchange:
			reply := Message{}
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
			request.Reply <- reply
		}
	}
}
