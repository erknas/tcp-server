package main

import (
	"log"
	"log/slog"
	"net"
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
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 5),
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	log.Println("PORT", s.listenAddr)

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

		log.Println("accepted connection:", conn.RemoteAddr().String())

		go s.read(conn)
	}
}

func (s *Server) read(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			slog.Info("reading error:", "err", err.Error())
			continue
		}

		msg := Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}

		s.msgch <- msg
	}
}
