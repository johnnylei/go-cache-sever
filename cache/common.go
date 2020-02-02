package cache

type Stat struct {
	Count int64 `json:"count"`
	KeySize int64 `json:"key_size"`
	ValueSize int64 `json:"value_size"`
}

func (_self *Stat)add(key string, value []byte)  {
	_self.Count++
	_self.KeySize += int64(len(key))
	_self.ValueSize += int64(len(value))
}

func (_self *Stat)del(key string, value []byte)  {
	_self.Count--
	_self.KeySize -= int64(len(key))
	_self.ValueSize -= int64(len(value))
}
