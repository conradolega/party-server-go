package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/textproto"
	"os"
	"time"

	"github.com/op/go-logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var logger = logging.MustGetLogger("logger")
var logFormat = logging.MustStringFormatter(
	`%{time} %{shortfunc} %{level} %{id:03x} %{message}`,
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
	messagesSent.With(prometheus.Labels{"type": "hello"}).Inc()
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
		messagesSent.With(prometheus.Labels{"type": msg}).Inc()
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
	messagesSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "messages_total",
			Help: "Number of messages sent",
		},
		[]string{"type"},
	)
)

func init() {
	prometheus.MustRegister(connections)
	prometheus.MustRegister(messagesSent)
}

func main() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, logFormat)
	logging.SetBackend(backendFormatter)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal(http.ListenAndServe(":3124", nil))
	}()

	logger.Info("Starting server...")
	server := Server{}
	server.Run()
}
