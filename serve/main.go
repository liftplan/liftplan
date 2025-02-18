package main

import (
	"embed"
	"log"
	"net/http"

	// _ "net/http/pprof"

	"github.com/liftplan/liftplan/serve/handler"
)

//go:embed static/*
var staticAssets embed.FS

func main() {
	http.HandleFunc("/", handler.Root())
	http.HandleFunc("/plan", handler.Plan())
	http.Handle("/static/", http.FileServerFS(staticAssets))
	s := http.Server{
		Addr: "0.0.0.0:9000",
	}
	log.Println(s.ListenAndServe())
}
