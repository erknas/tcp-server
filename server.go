package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitCh     chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	s.ln = ln
	go s.accept()

	<-s.quitCh

	return nil
}

func (s *Server) accept() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("acception error:", "err", err.Error())
			continue
		}

		log.Println("accepted connection:", conn.RemoteAddr())

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

		msg := buf[:n]

		fmt.Println(string(msg))
	}
}
