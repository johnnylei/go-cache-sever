package cache

import (
	"errors"
	"strings"
	"sync"
)

type Cache struct {
	*Stat
	Data map[string][]byte
	Lock sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		&Stat{},
		make(map[string][]byte),
		sync.RWMutex{},
	}
}

func (_self *Cache)Set(key string, value []byte) error  {
	key = strings.Trim(key, " ")
	if len(key) == 0 {
		return errors.New("key should not be empty")
	}

	_self.Lock.Lock()
	defer _self.Lock.Unlock()
	if _value, okay := _self.Data[key]; okay {
		_self.del(key, _value)
	}

	_self.Data[key] = value
	_self.add(key, value)
	return nil
}

func (_self *Cache)Get(key string) ([]byte, error)  {
	key = strings.Trim(key, " ")
	if len(key) == 0 {
		return nil, errors.New("key should not be empty")
	}

	_self.Lock.RLock()
	defer _self.Lock.RUnlock()
	return _self.Data[key], nil
}

func (_self *Cache)Del(key string) error  {
	key = strings.Trim(key, " ")
	if len(key) == 0 {
		return errors.New("key should not be empty")
	}

	_self.Lock.Lock()
	defer _self.Lock.Unlock()
	if value, okay := _self.Data[key]; okay {
		_self.del(key, value)
		delete(_self.Data, key)
	}

	return nil
}
