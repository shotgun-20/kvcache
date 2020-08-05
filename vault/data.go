package vault

import "time"

// Message - структура данных для двустороннего обмена с хранилищем,
// запрос/ответ
type Message struct {
	Key    string
	Value  interface{}
	Action string // SET, GET, DEL, POP
	Error  bool
	Reply  chan Message
}

// node - узел данных для очереди устаревания
type node struct {
	expire time.Time   // Окончание срока жизни
	Key    string      // Ключ
	Value  interface{} // Хранимое значение
	Prev   *node       // Предыдущий узел
	Next   *node       // Следующий узел
}

// Store - корневая структура для хранения данных
type Store struct {
	TTL      uint64           // Публично доступный TTL по умолчанию для инициализации, сек
	exchange chan Message     // Небуферизованный канал для синхронизации доступа
	ttl      time.Duration    // Время жизни узла по умолчанию, сек
	head     *node            // Голова, выходит первым
	tail     *node            // Добавлен последним
	flat     map[string]*node // Карта для быстрого доступа к значениям
}
