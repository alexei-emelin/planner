package api

import (
    "net/http"
    "planner/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    if id == "" {
        writeJson(w, map[string]string{"error": "не указан идентификатор"})
        return
    }

    if err := db.DeleteTask(id); err != nil {
        writeJson(w, map[string]string{"error": err.Error()})
        return
    }

    writeJson(w, map[string]string{})
}
