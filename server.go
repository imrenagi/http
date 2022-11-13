package http

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func Listen(address string) error {
	handler := func(conn net.Conn, r *Request) {
		res := "<h1>hello world</h1>"
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
			"Content-Length: %d\r\n"+
			"Content-Type: text/html\r\n"+
			"\r\n"+
			"%s", len(res), res)))
	}
	s := Server{Handler: handler}
	return s.Listen(address)
}

type Server struct {
	Handler func(conn net.Conn, r *Request)
}

func (s Server) Listen(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("can't start listener")
		return err
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("can't accept connection")
			return err
		}

		go s.handleConn(conn)
	}
}

func (s Server) handleConn(conn net.Conn) {
	request, err := ReadRequest(bufio.NewReader(conn))
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	s.Handler(conn, request)
}
