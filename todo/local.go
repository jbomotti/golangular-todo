// +build !appengine

// Create webserver for running Todo app
package main

import (
	"net/http"

	"github.com/jbomotti/golangular-todo/server"
)

func main() {
	server.RegisterHandlers()
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":8000", nil)
}
