package vault

import (
	"errors"
)

// init - инициализация хранилища с заданным TTL
func (store *Store) init() {
	if store.isInit == true {
		return
	}
	if store.TTL < 1 {
		store.TTL = 30
	}

	store.ttl = store.TTL
	store.exchange = make(chan Message)
	store.flat = make(map[string]*node)
	store.isInit = true

	// Создаём достаточное количество узлов-филлеров, чтобы механизм
	// устаревания не начал уничтожать записи преждевременно.
	var i uint64
	for i = 0; i < store.ttl; i++ {
		store.addNode("", "", false)
	}

	go store.control() // запускаем управление хранилищем
}

// addNode - добавляем новый узел в хвост очереди.
func (store *Store) addNode(key, value string, kind bool) error {
	var prev *node
	if store.tail != nil {
		prev = store.tail
	}
	store.tail = &node{Key: key, Value: value, Prev: prev, Kind: kind}
	if store.head == nil {
		store.head = store.tail
	}
	if prev != nil {
		prev.Next = store.tail
	}
	if kind == true {
		store.flat[key] = store.tail
	}
	return nil
}

// setNode - собираемся добавить новый узел.
// Если ключа ещё нет - просто добавляем в хвост.
// Если ключ есть - переносим в хвост.
func (store *Store) setNode(key, value string) error {
	if _, ok := store.flat[key]; ok == false {
		store.addNode(key, value, true)
		return nil
	}
	store.delNode(key)              // Сильно нерационально, удаление/создание объекта
	store.addNode(key, value, true) // Зато резко упрощает код
	return nil
}

// popNode - извлекаем узел из гловы очереди и уничтожаем его.
// Если Kind == 0 - то возвращаем ошибку.
// При ошибке устаревание приостанавливается на секунду.
func (store *Store) popNode() error {
	var err error
	if store.head == nil {
		return errors.New("No nodes")
	}
	if store.head.Kind == false {
		err = errors.New("Decay timeout")
	}
	if store.head.Key != "" {
		delete(store.flat, store.head.Key)
	}
	if store.head.Next != nil {
		store.head = store.head.Next
		store.head.Prev = nil
	}
	return err
}

// getNode - получить значение ключа, или ошибку, если такого нет.
func (store *Store) getNode(key string) (string, error) {
	if v, ok := store.flat[key]; ok != false {
		return v.Value, nil
	}
	return "", errors.New("No such key")
}

// delNode - удалить узел с указанным ключом.
// Если такого ключа не было - вернуть ошибку.
func (store *Store) delNode(key string) error {
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
