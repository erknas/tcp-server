package main

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message

	mu    sync.Mutex
	peers map[string]bool
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 5),
		peers:      make(map[string]bool),
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	slog.Info("starting server on", "PORT", s.listenAddr)

	s.ln = ln
	go s.accept()

	<-s.quitch
	close(s.msgch)

	return nil
}

func (s *Server) accept() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("acception error:", "err", err.Error())
			continue
		}

		slog.Info("accepted connection from", "addr", conn.RemoteAddr().String())

		s.mu.Lock()
		s.peers[conn.RemoteAddr().String()] = true
		fmt.Printf("accepted connections %v\n", s.peers)
		s.mu.Unlock()

		go s.read(conn)
	}
}

func (s *Server) read(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			slog.Error("reading error:", "err", err.Error())
			continue
		}

		msg := Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}

		s.msgch <- msg
	}
}
