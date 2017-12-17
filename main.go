package main

import (
	"fmt"
	"net"
)

type Server struct {
	clients []net.Conn
}

func (s Server) Run() {
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
		s.clients = append(s.clients, conn)
		go s.Handle(conn)
	}
}

func (s Server) Handle(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("Hello\n"))
}

func main() {
	fmt.Println("Starting server...")
	server := Server{}
	server.Run()
}
