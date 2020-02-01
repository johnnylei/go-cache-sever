package main

func main() {
	server := NewTcpServer()
	server.Run("0.0.0.0:12345")
	//listener, _ := net.Listen("tcp", "0.0.0.0:12345")
	//client, _ := listener.Accept()
	//defer client.Close()
	//defer listener.Close()
	//buffer := make([]byte, 1024)
	//client.Read(buffer)
	//fmt.Println(string(buffer))
	//client.Write([]byte("succeed"))
}
