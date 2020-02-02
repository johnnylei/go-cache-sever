package main

import "flag"

func main() {
	cacheType := flag.String("t", "default", "cache type")
	flag.Parse()
	server := NewTcpServer(*cacheType)
	server.Run("0.0.0.0:12345")
}
