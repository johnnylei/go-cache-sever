package main

import (
	"bufio"
	"github.com/johnnylei/go-cache-sever/cache"
	"io"
	"net"
	"strconv"
	"strings"
)

func NewTcpServer(t string) *TcpServer  {
	m := map[string]func()cache.CacheInterface{
		"default": func() cache.CacheInterface {
			return cache.NewCache()
		},
		"rocksdb": func() cache.CacheInterface {
			return cache.NewRocksdbCache()
		},
	}
	return &TcpServer{
		m[t](),
	}
}

type TcpServer struct {
	cache cache.CacheInterface
}

func (_self *TcpServer)Run(address string) error  {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			return err
		}

		go _self.process(connection)
	}
}

func (_self *TcpServer) process(conn net.Conn)  {
	defer func() {
		err := recover()
		if err != nil {
			errMessage, _ := err.(string)
			writer := bufio.NewWriter(conn)
			writer.Write([]byte(errMessage))
		}

		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	op, err := reader.ReadByte()
	if err != nil {
		panic(err)
	}

	switch op {
	case 'S':
		if err := _self.set(reader); err != nil {
			panic(err)
		}

		writer.Write([]byte("succeed"))
		break
	case 'G':
		value, err := _self.get(reader)
		if err != nil {
			panic(err)
		}

		writer.Write(value)
		break
	case 'D':
		if err := _self.del(reader); err != nil {
			panic(err)
		}

		writer.Write([]byte("succeed"))
		break
	default:
		if _, err := writer.WriteString("unknow operator"); err != nil {
			panic(err)
		}
		break
	}
	writer.Flush()
}

func (_self *TcpServer) get(reader *bufio.Reader) ([]byte, error)  {
	keyLength,  err:= _self.readLength(reader)
	if err != nil {
		return nil, err
	}

	keyBuffer, err := _self.read(reader, keyLength)
	if err != nil {
		return nil, err
	}

	valueBuffer, err := _self.cache.Get(string(keyBuffer))
	if err != nil {
		return nil, err
	}

	return valueBuffer, nil
}

func (_self *TcpServer) set(reader *bufio.Reader) error  {
	keyLength,  err:= _self.readLength(reader)
	if err != nil {
		return err
	}

	valueLength, err := _self.readLength(reader)
	if err != nil {
		return err
	}

	keyBuffer, err := _self.read(reader, keyLength)
	if err != nil {
		return err
	}

	valueBuffer, err := _self.read(reader, valueLength)
	if err != nil {
		return err
	}

	return _self.cache.Set(string(keyBuffer), valueBuffer)
}

func (_self *TcpServer) del(reader *bufio.Reader) error  {
	keyLength,  err:= _self.readLength(reader)
	if err != nil {
		return err
	}

	keyBuffer, err := _self.read(reader, keyLength)
	if err != nil {
		return err
	}

	return _self.cache.Del(string(keyBuffer))
}

func (_self *TcpServer) read(reader *bufio.Reader, keyLen int) ([]byte, error)  {
	buffer := make([]byte, keyLen)
	_, err := io.ReadFull(reader, buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func (_self *TcpServer) readLength(reader *bufio.Reader) (int, error) {
	lenStr, err := reader.ReadString(' ')
	if err != nil {
		return 0, err
	}

	length, err := strconv.Atoi(strings.TrimSpace(lenStr))
	if err != nil {
		return 0, err
	}

	return length, nil
}
