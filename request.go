package http

import (
	"bufio"
	"fmt"
	"net/textproto"
	"strings"
)

func ReadRequest(buf *bufio.Reader) (*Request, error) {
	tp := textproto.NewReader(buf)
	line, err := tp.ReadLine()
	if err != nil {
		return nil, err
	}

	sl := strings.Split(line, " ")
	if len(sl) < 3 {
		return nil, fmt.Errorf("invalid http data")
	}
	method, path, protoVersion := sl[0], sl[1], sl[2]

	return &Request{
		Method: method,
		URI:    path,
		Proto:  protoVersion,
	}, nil
}

type Request struct {
	Method string
	URI    string
	Proto  string
}
