package cache

// #include "rocksdb/c.h"
// #include <stdlib.h>
// #cgo CFLAGS: -I/root/programs/rocksdb-5.11.2/include
// #cgo LDFLAGS: -L/root/programs/rocksdb-5.11.2 -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"errors"
	"regexp"
	"runtime"
	"strconv"
	"unsafe"
)

type rocksdbCache struct {
	db *C.rocksdb_t
	ro *C.rocksdb_readoptions_t
	wo *C.rocksdb_writeoptions_t
	e *C.char
}

func NewRocksdbCache() *rocksdbCache {
	options := C.rocksdb_options_create()
	C.rocksdb_options_increase_parallelism(options, C.int(runtime.NumCPU()))
	C.rocksdb_options_set_create_if_missing(options, 1)
	var e *C.char
	db := C.rocksdb_open(options, C.CString("/mnt/rocksdb"), &e)
	if e != nil {
		panic(C.GoString(e))
	}

	C.rocksdb_options_destroy(options)
	return &rocksdbCache{
		db,
		C.rocksdb_readoptions_create(),
		C.rocksdb_writeoptions_create(),
		e,
	}
}

func (_self *rocksdbCache)Get(key string) ([]byte, error)  {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	var length C.size_t
	v := C.rocksdb_get(_self.db, _self.ro, k, C.size_t(len(key)), &length, &_self.e)
	if _self.e != nil {
		return nil, errors.New(C.GoString(_self.e))
	}
	defer C.free(unsafe.Pointer(v))
	return C.GoBytes(unsafe.Pointer(v), C.int(length)), nil
}

func (_self *rocksdbCache)Del(key string) error  {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))

	C.rocksdb_delete(_self.db, _self.wo, k, C.size_t(len(key)), &_self.e)
	if _self.e != nil {
		return errors.New(C.GoString(_self.e))
	}
	return nil
}

func (_self *rocksdbCache)Set(key string, value []byte) error  {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	v := C.CBytes(value)
	defer C.free(unsafe.Pointer(v))

	C.rocksdb_put(_self.db, _self.wo, k, C.size_t(len(key)), (*C.char)(v), C.size_t(len(value)), &_self.e)
	if _self.e != nil {
		return errors.New(C.GoString(_self.e))
	}
	return nil
}

func (c *rocksdbCache) GetStat() Stat {
	k := C.CString("rocksdb.aggregated-table-properties")
	defer C.free(unsafe.Pointer(k))
	v := C.rocksdb_property_value(c.db, k)
	defer C.free(unsafe.Pointer(v))
	p := C.GoString(v)
	r := regexp.MustCompile(`([^;]+)=([^;]+);`)
	s := Stat{}
	for _, submatches := range r.FindAllStringSubmatch(p, -1) {
		if submatches[1] == " # entries" {
			s.Count, _ = strconv.ParseInt(submatches[2], 10, 64)
		} else if submatches[1] == " raw key size" {
			s.KeySize, _ = strconv.ParseInt(submatches[2], 10, 64)
		} else if submatches[1] == " raw value size" {
			s.ValueSize, _ = strconv.ParseInt(submatches[2], 10, 64)
		}
	}
	return s
}
