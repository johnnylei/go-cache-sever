package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
)

func buildCommand() (string, error)  {
	op := flag.String("c", "get", "operator")
	key := flag.String("k", "", "key")
	value := flag.String("v", "", "value")
	flag.Parse()
	operatorMap := map[string]func(string, string)string{
		"get" : func(key string, value string) string {
			kenLength := len(key)
			return fmt.Sprintf("G%d %s", kenLength, key)
		},
		"set" : func(key string, value string) string {
			kenLength := len(key)
			valueLength := len(value)
			return fmt.Sprintf("S%d %d %s%s", kenLength, valueLength, key, value)
		},
		"del" : func(key string, value string) string {
			kenLength := len(key)
			return fmt.Sprintf("D%d %s", kenLength, key)
		},
	}
	if _, okay := operatorMap[*op]; !okay {
		return "", errors.New("operator should be in array get, set, del")
	}
	return operatorMap[*op](*key, *value), nil
}

func main() {
	command, err := buildCommand()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(command)

	client, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	writer := bufio.NewWriter(client)
	if _, err = writer.WriteString(command); err != nil {
		fmt.Println(err.Error())
		return
	}
	writer.Flush()


	reader := bufio.NewReader(client)
	buffer := make([]byte, 1024)
	if _, err := reader.Read(buffer); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf(string(buffer))
}
