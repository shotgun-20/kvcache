package vault

// Message - структура данных для двустороннего обмена с хранилищем,
// запрос/ответ
type Message struct {
	Key    string
	Value  string
	Action string // SET, GET, DEL, POP
	Error  bool
	Reply  chan Message
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
	Exchange chan Message     // Небуферизованный канал для синхронизации доступа
	ttl      uint64           // Время жизни узла, сек
	head     *node            // Голова, выходит первым
	tail     *node            // Добавлен последним
	flat     map[string]*node // Карта для быстрого доступа к значениям
}
