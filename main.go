package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Starting server...")
	listener, err := net.Listen("tcp", "0.0.0.0:3123")
	if err != nil {
		fmt.Println(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Client %v connected.\n", conn.RemoteAddr())
		conn.Write([]byte("Hello\n"))
		conn.Close()
	}
}
