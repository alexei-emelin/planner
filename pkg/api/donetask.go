package api

import (
    "net/http"
    "planner/pkg/db"
	"time"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
        return
    }

    id := r.URL.Query().Get("id")
    if id == "" {
        writeJson(w, map[string]string{"error": "не указан идентификатор"})
        return
    }

    task, err := db.GetTask(id)
    if err != nil {
        writeJson(w, map[string]string{"error": err.Error()})
        return
    }

    // Если правило повторения пустое — просто удаляем задачу
    if task.Repeat == "" {
        if err := db.DeleteTask(id); err != nil {
            writeJson(w, map[string]string{"error": err.Error()})
            return
        }
        writeJson(w, map[string]string{})
        return
    }

    // Для периодических — вычисляем следующую дату
    now := time.Now()
    next, err := NextDate(now, task.Date, task.Repeat)
    if err != nil {
        writeJson(w, map[string]string{"error": err.Error()})
        return
    }

    if err := db.UpdateDate(next, id); err != nil {
        writeJson(w, map[string]string{"error": err.Error()})
        return
    }

    writeJson(w, map[string]string{})
}
