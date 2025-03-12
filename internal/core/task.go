package core

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
)

// Task представляет задачу.
type Task struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	AuthorID   int    `json:"authorID"`
	AssignedID int    `json:"assignedID"`
	Opened     int64  `json:"opened"`
	Closed     int64  `json:"closed"`
}

// TaskRepository определяет методы для работы с хранилищем задач.
type TaskRepository interface {
	GetTasks(taskID, authorID int) ([]Task, error)
	CreateTask(task Task) (int, error)
	CreateTasks(tasks []Task) ([]int, error)
	UpdateTask(task Task) error
	DeleteTask(taskID int) error
}

// TaskService определяет методы бизнес-логики задач.
type TaskService interface {
	GetTasks(taskID, authorID int) ([]Task, error)
	CreateTask(task Task) (int, error)
	CreateTasks(tasks []Task) ([]int, error)
	UpdateTask(task Task) error
	DeleteTask(taskID int) error
}

type taskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) GetTasks(taskID, authorID int) ([]Task, error) {
	return s.repo.GetTasks(taskID, authorID)
}

func (s *taskService) CreateTask(task Task) (int, error) {
	if task.Title == "" {
		return 0, ErrInvalidInput
	}
	return s.repo.CreateTask(task)
}

func (s *taskService) CreateTasks(tasks []Task) ([]int, error) {
	if len(tasks) == 0 {
		return nil, ErrInvalidInput
	}
	return s.repo.CreateTasks(tasks)
}

func (s *taskService) UpdateTask(task Task) error {
	if task.ID == 0 || task.Title == "" {
		return ErrInvalidInput
	}
	return s.repo.UpdateTask(task)
}

func (s *taskService) DeleteTask(taskID int) error {
	if taskID == 0 {
		return ErrInvalidInput
	}
	return s.repo.DeleteTask(taskID)
}
