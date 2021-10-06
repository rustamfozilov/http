package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type HandlerFunction func(conn net.Conn)

type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandlerFunction
}

func NewServer(addr string) *Server {
	return &Server{
		addr:     addr,
		mu:       sync.RWMutex{},
		handlers: make(map[string]HandlerFunction),
	}
}

func (s *Server) Register(path string, handler HandlerFunction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer func() {
		cerr := listener.Close()
		if cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Println(err)
		}

	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		path, err := getPath(conn)
		handle, ok := s.handlers[path]
		if !ok {
			return err
		}
		 go handle(conn)
	}

	return nil
}

func getPath(conn net.Conn) (path string, err error) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Println(err)
		}
	}()
	buf := make([]byte, 4096)

	read, err := conn.Read(buf)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return "", err
	}
	data := buf[:read]
	requestLineDelim := []byte{'\r', '\n'}

	requestLineEnd := bytes.Index(data, requestLineDelim)
	if requestLineEnd == -1 {
		return "", err
	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return"", err
	}
	path = parts[1]
	return path, nil
}
