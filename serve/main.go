package main

import (
	"embed"
	"net/http"

	// _ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/liftplan/liftplan/serve/handler"
)

//go:embed static/*
var staticAssets embed.FS

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.HandleFunc("/", handler.Root())
	r.HandleFunc("/v2", handler.RootV2())
	r.HandleFunc("/plan", handler.Plan())
	r.Handle("/static/*", http.FileServerFS(staticAssets))
	http.ListenAndServe(":9000", r)
}
