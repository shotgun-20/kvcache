package vault

import "errors"

// addNode - добавляем новый узел в хвост очереди.
func (store *store) addNode(key, value string, kind bool) {
	prev := store.tail
	store.tail = &node{Key: key, Value: value, Prev: prev, Kind: kind}
	prev.Next = store.tail
	if kind == true {
		store.flat[key] = store.tail
	}
}

// setNode - собираемся добавить новый узел.
// Если ключа ещё нет - просто добавляем в хвост.
// Если ключ есть - переносим в хвост.
func (store *store) setNode(key, value string) error {
	if _, ok := store.flat[key]; ok == false {
		store.addNode(key, value, true)
		return nil
	}
	// Вместо того, чтобы удалить и пересоздать, просто правим ссылки.
	// Так будет эффективней, чем удалять и создавать объекты данных.
	node := store.flat[key]
	head := store.head
	prev := node.Prev
	next := node.Next
	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
	}
	if store.head == node {
		store.head = prev
	}
	if store.tail == node {
		store.tail = next
	}
	store.head = node
	node.Next = head
	return nil
}

// popNode - извлекаем узел из гловы очереди и уничтожаем его.
// Если Kind == 0 - то возвращаем ошибку.
// При ошибке устаревание приостанавливается на секунду.
func (store *store) popNode() error {
	if store.head == nil {
		return errors.New("No nodes")
	}
	if store.head.Kind == false {
		return errors.New("Decay timeout")
	}
	key := store.head.Key
	delete(store.flat, key)
	store.head = store.head.Next
	return nil
}

// getNode - получить значение ключа, или ошибку, если такого нет.
func (store *store) getNode(key string) (string, error) {
	if v, ok := store.flat[key]; ok != false {
		return v.Value, nil
	}
	return "", errors.New("No such key")
}

// delNode - удалить узел с указанным ключом.
// Если такого ключа не было - вернуть ошибку.
func (store *store) delNode(key string) error {
	if _, ok := store.flat[key]; ok == false {
		return errors.New("No such key")
	}
	node := store.flat[key]
	prev := node.Prev
	next := node.Next
	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
	}
	if store.head == node {
		store.head = prev
	}
	if store.tail == node {
		store.tail = next
	}
	delete(store.flat, key)
	node = nil
	return nil
}
