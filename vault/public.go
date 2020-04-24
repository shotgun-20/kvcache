/* Методы, доступные внешним сущностям. */

package vault

import (
	"errors"
)

// GetValue - получить значение ключа из хранилища
func (store *Store) GetValue(key string) (interface{}, error) {
	if store.isInit == false {
		store.init()
	}
	reply := make(chan Message)
	msg := Message{Action: "GET", Reply: reply, Key: key}
	store.exchange <- msg
	resp := <-reply
	if resp.Error != true {
		return resp.Value, nil
	}
	return "", errors.New("No such key")
}

// SetValue - установить/обновить значение ключа
func (store *Store) SetValue(key string, value interface{}) error {
	if store.isInit == false {
		store.init()
	}
	reply := make(chan Message)
	msg := Message{Action: "SET", Reply: reply, Key: key, Value: value}
	store.exchange <- msg
	resp := <-reply
	if resp.Error != true {
		return nil
	}
	return errors.New("Cannot set value")
}

// DelValue - удалить ключ из хранилища
func (store *Store) DelValue(key string) error {
	if store.isInit == false {
		store.init()
	}
	reply := make(chan Message)
	msg := Message{Action: "DEL", Reply: reply, Key: key}
	store.exchange <- msg
	resp := <-reply
	if resp.Error != true {
		return nil
	}
	return errors.New("There was no such key")
}
