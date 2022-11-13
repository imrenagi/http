package http

import (
	"bufio"
	"fmt"
	"strings"
)

func ReadRequest(buf *bufio.Reader) (*Request, error) {
	// parsing the first line of http request
	var fl []byte
	for {
		line, more, err := buf.ReadLine()
		if err != nil {
			return nil, err
		}
		if line == nil && !more {
			fl = line
			break
		}
		fl = append(fl, line...)
		if !more {
			break
		}
	}

	sl := strings.Split(string(fl), " ")
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
