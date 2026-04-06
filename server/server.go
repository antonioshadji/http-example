package server

import (
	"fmt"
	"net/http"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<!DOCTYPE html><html><head><title>http-example</title></head><body><h1>I'm running on your machine</h1></body></html>")
	})
	return mux
}
