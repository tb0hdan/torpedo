package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	listen = flag.String("listen", "localhost:8000", "HTTP Server listen address")
	directory = flag.String("directory", "./", "Serve this directory")
)

func main() {
	flag.Parse()
	log.Fatal(http.ListenAndServe(*listen, http.FileServer(http.Dir(*directory))))
}
