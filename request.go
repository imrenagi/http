package http

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ReadRequest(buf *bufio.Reader) (*Request, error) {
	// parsing the first line of http request
	fl, err := readLine(buf)
	if err != nil {
		return nil, err
	}

	sl := strings.Split(string(fl), " ")
	if len(sl) < 3 {
		return nil, fmt.Errorf("invalid http data")
	}
	method, path, protoVersion := sl[0], sl[1], sl[2]

	headers, err := readHeader(buf)
	if err != nil {
		return nil, err
	}

	// get content length
	length, _ := strconv.Atoi(headers["Content-Length"][0])
	body := make([]byte, length)

	_, err = io.ReadFull(buf, body)
	if err != nil {
		return nil, err
	}

	return &Request{
		Method:  method,
		URI:     path,
		Proto:   protoVersion,
		Headers: headers,
		Body:    io.NopCloser(bytes.NewReader(body)),
	}, nil
}

func readHeader(buf *bufio.Reader) (Headers, error) {
	headers := make(map[string][]string)
	var foundEmptyLine bool
	for !foundEmptyLine {
		l, err := readLine(buf)
		if err != nil {
			return nil, err
		}

		if string(l) == "" {
			foundEmptyLine = true
			break
		}

		k, v, ok := bytes.Cut(l, []byte(":"))
		if !ok {
			return nil, fmt.Errorf("invalid header format")
		}

		vs := strings.TrimLeft(string(v), " \t")

		val, ok := headers[string(k)]
		if !ok {
			headers[string(k)] = []string{vs}
		} else {
			values := append(val, vs)
			headers[string(k)] = values
		}

	}
	return headers, nil
}

func readLine(buf *bufio.Reader) ([]byte, error) {
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
	return fl, nil
}

type Headers map[string][]string

type Request struct {
	Method  string
	Headers Headers
	URI     string
	Proto   string
	Body    io.ReadCloser
}
