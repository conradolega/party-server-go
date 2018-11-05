package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/textproto"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		connections.With(prometheus.Labels{}).Inc()
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

var (
	connections = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "connections_total",
			Help: "Number of clients connected.",
		},
		[]string{},
	)
)

func init() {
	prometheus.MustRegister(connections)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal(http.ListenAndServe(":3124", nil))
	}()

	fmt.Println("Starting server...")
	server := Server{}
	server.Run()
}
