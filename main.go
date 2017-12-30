package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"time"
)

type Server struct {
	clients []net.Conn
}

func (s Server) Run() {
	listener, err := net.Listen("tcp", "0.0.0.0:3123")
	if err != nil {
		fmt.Println(err)
	}

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for range ticker.C {
			s.SendToAll("ALCOHOL\n")
		}
	}()

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
	conn.Write([]byte("Hello\n"))
	s.SendToAll(fmt.Sprintf("%v", conn) + " has connected\n")

	reader := textproto.NewReader(bufio.NewReader(conn))
	for {
		msg, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(msg)
	}
}

func (s Server) SendToAll(msg string) {
	for _, client := range s.clients {
		client.Write([]byte(msg))
	}
}

func main() {
	fmt.Println("Starting server...")
	server := Server{}
	server.Run()
}
