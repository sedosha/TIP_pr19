package http

import (
    "encoding/json"
    "net/http"
    "sync"
    "time"
    
    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "tech-ip-sem2/services/tasks/internal/service"
    "tech-ip-sem2/shared/middleware"
)

type TaskHandler struct {
    mu    sync.RWMutex
    tasks map[string]service.Task
}

func NewTaskHandler() *TaskHandler {
    return &TaskHandler{
        tasks: make(map[string]service.Task),
    }
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    requestID := middleware.GetRequestID(r.Context())
    
    var req service.CreateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid request format"}`, http.StatusBadRequest)
        return
    }
    
    if req.Title == "" {
        http.Error(w, `{"error":"title is required"}`, http.StatusBadRequest)
        return
    }
    
    h.mu.Lock()
    defer h.mu.Unlock()
    
    id := uuid.New().String()[:8]
    task := service.Task{
        ID:          id,
        Title:       req.Title,
        Description: req.Description,
        DueDate:     req.DueDate,
        Done:        false,
        CreatedAt:   time.Now(),
    }
    
    h.tasks[id] = task
    
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Request-ID", requestID)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
    requestID := middleware.GetRequestID(r.Context())
    
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    tasks := make([]service.Task, 0, len(h.tasks))
    for _, task := range h.tasks {
        tasks = append(tasks, task)
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Request-ID", requestID)
    json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
    requestID := middleware.GetRequestID(r.Context())
    id := chi.URLParam(r, "id")
    
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    task, exists := h.tasks[id]
    if !exists {
        http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Request-ID", requestID)
    json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    requestID := middleware.GetRequestID(r.Context())
    id := chi.URLParam(r, "id")
    
    var req service.UpdateTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid request format"}`, http.StatusBadRequest)
        return
    }
    
    h.mu.Lock()
    defer h.mu.Unlock()
    
    task, exists := h.tasks[id]
    if !exists {
        http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
        return
    }
    
    if req.Title != nil {
        task.Title = *req.Title
    }
    if req.Description != nil {
        task.Description = *req.Description
    }
    if req.Done != nil {
        task.Done = *req.Done
    }
    
    h.tasks[id] = task
    
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Request-ID", requestID)
    json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    requestID := middleware.GetRequestID(r.Context())
    id := chi.URLParam(r, "id")
    
    h.mu.Lock()
    defer h.mu.Unlock()
    
    if _, exists := h.tasks[id]; !exists {
        http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
        return
    }
    
    delete(h.tasks, id)
    
    w.Header().Set("X-Request-ID", requestID)
    w.WriteHeader(http.StatusNoContent)
}
