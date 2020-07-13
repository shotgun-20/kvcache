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
	Prev   *node       // Ближе к голове, извлекается раньше
	Next   *node       // Ближе к хвосту, выйдёт позже
}

// Store - корневая структура для хранения данных
type Store struct {
	TTL      uint64           // Публично доступный TTL по умолчанию для инициализации
	exchange chan Message     // Небуферизованный канал для синхронизации доступа
	ttl      uint64           // Время жизни узла по умолчанию, сек
	head     *node            // Голова, выходит первым
	tail     *node            // Добавлен последним
	flat     map[string]*node // Карта для быстрого доступа к значениям
}
