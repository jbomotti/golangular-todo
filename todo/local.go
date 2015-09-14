// +build !appengine

// Create webserver for running Todo app
package main

import (
	"net/http"

	"github.com/jbomotti/todo-test/server"
)

func main() {
	server.RegisterHandlers()
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":8000", nil)
}
