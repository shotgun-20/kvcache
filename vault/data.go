package vault

import "time"

// Message - структура данных для двустороннего обмена с хранилищем,
// запрос/ответ
type Message struct {
	Key    string
	Value  string
	Action string
	Code   int
}

// node - узел данных для очереди устаревания
type node struct {
	Kind  bool   // 0 - таймаут, 1 - данные
	Value string // Хранимое значение
}

// Store - корневая структура для хранения данных
type Store struct {
	Exchange chan Message // Небуферизованный канал для синхронизации доступа
	TTL      time.Time    // Время жизни узла
	Head     *node
	Tail     *node
	Map      map[string]node // Карта для быстрого доступа к значениям
}
