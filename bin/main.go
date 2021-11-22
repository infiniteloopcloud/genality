package main

import (
	"flag"
	"log"

	"github.com/infiniteloopcloud/genality/bin/http"
)

type Communicator interface {
	Serve() error
}

var (
	communicator string
)

func init() {
	flag.StringVar(&communicator, "communicator", "http", "Set the type of communicator (http)")
	flag.Parse()
}

func main() {
	var c Communicator
	switch communicator {
	case "http":
		c = http.Communicator{}
	}
	if c != nil {
		if err := c.Serve(); err != nil {
			log.Fatal(err)
		}
	}
	log.Fatal("Serve failed")
}
