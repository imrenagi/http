package main

import (
	"log"

	"github.com/imrenagi/http"
)

func main() {
	err := http.Listen("localhost:8080")
	if err != nil {
		log.Fatalf(err.Error())
	}
}
