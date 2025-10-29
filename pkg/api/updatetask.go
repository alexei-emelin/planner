package api

import (
    "encoding/json"
    "planner/pkg/db"
    "net/http"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
    var task db.Task
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&task); err != nil {
        writeJson(w, map[string]string{"error": "ошибка десериализации JSON"}, http.StatusBadRequest)
        return
    }

    if task.Title == "" {
        writeJson(w, map[string]string{"error": "не указан заголовок задачи"}, http.StatusBadRequest)
        return
    }

    if err := checkDate(&task); err != nil {
        writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
        return
    }

    err := db.UpdateTask(&task)
    if err != nil {
        writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
        return
    }

    // успешное обновление
    writeJson(w, map[string]string{}, http.StatusBadRequest)
}
