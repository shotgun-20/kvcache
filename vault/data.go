package vault

import "time"

// Message - структура данных для двустороннего обмена с хранилищем,
// запрос/ответ
type Message struct {
	Key    string
	Value  string
	Action string
	Error  bool
}

// node - узел данных для очереди устаревания
type node struct {
	Kind  bool   // 0 - таймаут, 1 - данные
	Key   string // Ключ
	Value string // Хранимое значение
	Prev  *node  // Ближе к голове, извлекается раньше
	Next  *node  // Ближе к хвосту, выйдёт позже
}

// Store - корневая структура для хранения данных
type store struct {
	exchange chan Message  // Небуферизованный канал для синхронизации доступа
	ttl      time.Duration // Время жизни узла
	head     *node
	tail     *node
	flat     map[string]*node // Карта для быстрого доступа к значениям
}
