package main

import (
	"go-specs-greet/adapters/httpserver"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(httpserver.Handler)
	http.ListenAndServe(":8080", handler)
}
