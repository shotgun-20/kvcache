package vault

// setNode - добавляем новый узел в очередь.
// Если такой ключ уже есть - переносим узел в начало очереди.
func (store *Store) setNode(key, value string) error {
	return nil
}

// popNode - извлекаем узел из очереди и уничтожаем его.
// Если Kind == 0 - то возвращаем ошибку.
// При ошибке устаревание приостанавливается на секунду.
func (store *Store) popNode() error {
	return nil
}

// getNode - получить значение ключа, или ошибку, если такого нет.
func (store *Store) getNode(key string) (string, error) {
	return "", nil
}

// delNode - удалить узел с указанным ключом.
// Если такого ключа не было - вернуть ошибку.
func (store *Store) delNode(key string) error {
	return nil
}
