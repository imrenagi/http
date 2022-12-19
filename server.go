package http

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
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

	postHandler := func(conn net.Conn, r *Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		b64 := base64.StdEncoding.EncodeToString(data)
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
			"Content-Length: %d\r\n"+
			"Content-Type: text/html\r\n"+
			"\r\n"+
			"%s", len(b64), b64)))
	}

	s := Server{
		Handler:     handler,
		PostHandler: postHandler,
	}
	return s.Listen(address)
}

type Server struct {
	Handler     func(conn net.Conn, r *Request)
	PostHandler func(conn net.Conn, r *Request)
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
	switch request.Method {
	case "GET":
		s.Handler(conn, request)
	case "POST":
		s.PostHandler(conn, request)
	}
}
