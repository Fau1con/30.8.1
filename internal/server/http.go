package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"30.8.1/internal/core"
	"github.com/gorilla/mux"
)

type HTTPServer struct {
	router  *mux.Router
	service core.TaskService
}

func NewHTTPServer(service core.TaskService) *HTTPServer {
	s := &HTTPServer{
		router:  mux.NewRouter(),
		service: service,
	}
	s.routes()
	return s
}

func (s *HTTPServer) routes() {
	s.router.HandleFunc("/tasks", s.getTasks).Methods("GET")
	s.router.HandleFunc("/tasks", s.createTask).Methods("POST")
	s.router.HandleFunc("/tasks/batch", s.createTasks).Methods("POST")
	s.router.HandleFunc("/tasks/{id}", s.updateTask).Methods("PUT")
	s.router.HandleFunc("/tasks/{id}", s.deleteTask).Methods("DELETE")
}

func (s *HTTPServer) Router() http.Handler {
	return s.router
}

func (s *HTTPServer) getTasks(w http.ResponseWriter, r *http.Request) {
	taskID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	authorID, _ := strconv.Atoi(r.URL.Query().Get("authorID"))

	tasks, err := s.service.GetTasks(taskID, authorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func (s *HTTPServer) createTask(w http.ResponseWriter, r *http.Request) {
	var task core.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id, err := s.service.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (s *HTTPServer) createTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []core.Task
	if err := json.NewDecoder(r.Body).Decode(&tasks); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ids, err := s.service.CreateTasks(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string][]int{"ids": ids})
}

func (s *HTTPServer) updateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task core.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	task.ID = taskID

	if err := s.service.UpdateTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (s *HTTPServer) deleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := s.service.DeleteTask(taskID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
