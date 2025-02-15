package main

import (
	"log"
	"net/http"

	// _ "net/http/pprof"

	"github.com/liftplan/liftplan/serve/handler"
)

func main() {
	http.HandleFunc("/", handler.Root())
	http.HandleFunc("/plan", handler.Plan())
	s := http.Server{
		Addr: "0.0.0.0:9000",
	}
	log.Println(s.ListenAndServe())
}
