package api

import (
    "net/http"
    "planner/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        writeJson(w, map[string]string{"error": "не указан идентификатор"}, http.StatusBadRequest)
        return
    }

    task, err := db.GetTask(id)
    if err != nil {
        writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
        return
    }

    writeJson(w, task, http.StatusBadRequest)
}
