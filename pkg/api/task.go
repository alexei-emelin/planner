package api

import (
	"net/http"
	"encoding/json"
	"planner/pkg/db"
)
func taskHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        addTaskHandler(w, r)
    case http.MethodGet:
        getTaskHandler(w, r)
    case http.MethodPut:
        updateTaskHandler(w, r)
    case http.MethodDelete:
        deleteTaskHandler(w, r)
    default:
        http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
    }
}


// writeJson возвращает json-ответ
func writeJson(w http.ResponseWriter, data any) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    json.NewEncoder(w).Encode(data)
}

type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
        return
    }

    tasks, err := db.Tasks(50)
    if err != nil {
        writeJson(w, map[string]string{"error": err.Error()})
        return
    }

    writeJson(w, TasksResp{Tasks: tasks})
}