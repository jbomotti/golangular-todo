package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jbomotti/golangular-todo/todo/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/jbomotti/golangular-todo/todo/Godeps/_workspace/src/github.com/jbomotti/golangular-todo/task"
)

var tasks = task.NewList()

const PathPrefix = "/task/"

func RegisterHandlers() {
	r := mux.NewRouter()
	r.HandleFunc(PathPrefix, errorHandler(ListTasks)).Methods("GET")
	r.HandleFunc(PathPrefix, errorHandler(NewTask)).Methods("POST")
	r.HandleFunc(PathPrefix+"{id}", errorHandler(GetTask)).Methods("GET")
	r.HandleFunc(PathPrefix+"{id}", errorHandler(UpdateTask)).Methods("PUT")
	http.Handle(PathPrefix, r)
}

type badRequest struct{ error }
type notFound struct{ error }

func errorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch err.(type) {
		case badRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case notFound:
			http.Error(w, "Task not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(w, "Oops, something went wrong", http.StatusInternalServerError)
		}
	}
}

// Returns an object containing all Tasks when a GET request is sent to '/task/'
func ListTasks(w http.ResponseWriter, r *http.Request) error {
	res := struct{ Tasks []*task.Task }{tasks.All()}
	return json.NewEncoder(w).Encode(res)
}

// Creates a new Task when a POST request is sent to '/task/'
func NewTask(w http.ResponseWriter, r *http.Request) error {
	req := struct{ Title string }{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequest{err}
	}
	t, err := task.NewTask(req.Title)
	if err != nil {
		return badRequest{err}
	}
	return tasks.Save(t)
}

// Helper for getting ID variable from the request url, makes use of ParseInt from strconv package
func parseID(r *http.Request) (int64, error) {
	txt, ok := mux.Vars(r)["id"]
	if !ok {
		return 0, fmt.Errorf("Task ID not found")
	}
	return strconv.ParseInt(txt, 10, 0)
}

// Handles GET requests for a specific task '/task/{id}'
func GetTask(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	log.Println("Task is #", id)
	if err != nil {
		return badRequest{err}
	}
	t, ok := tasks.Find(id)
	log.Println("Found", ok)

	if !ok {
		return notFound{}
	}
	return json.NewEncoder(w).Encode(t)
}

// Handles PUT requests to a specific task '/task/{id}'
func UpdateTask(w http.ResponseWriter, r *http.Request) error {
	id, err := parseID(r)
	if err != nil {
		return badRequest{err}
	}
	var t task.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return badRequest{err}
	}
	if t.ID != id {
		return badRequest{fmt.Errorf("Task IDs do not match")}
	}
	if _, ok := tasks.Find(id); !ok {
		return notFound{}
	}
	return tasks.Save(&t)
}
