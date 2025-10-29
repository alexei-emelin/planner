package api

import (
    "encoding/json"
    "net/http"
    "planner/pkg/db"
    "time"
)

// Обработчик для POST /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
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

    id, err := db.AddTask(&task)
    if err != nil {
        writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
        return
    }
    writeJson(w, map[string]interface{}{"id": id}, http.StatusBadRequest)
}

// Проверка и коррекция даты задачи, валидация repeat и вычисление nextDate
func checkDate(task *db.Task) error {
    now := time.Now()
    if task.Date == "" {
        task.Date = now.Format(DateFormat)
    }
    t, err := time.Parse(DateFormat, task.Date)
    if err != nil {
        return err
    }
    var next string
    if task.Repeat != "" {
        var nextErr error
        next, nextErr = NextDate(now, task.Date, task.Repeat)
        if nextErr != nil {
            return nextErr
        }
    }
    if afterNow(now, t) {
        if len(task.Repeat) == 0 {
            task.Date = now.Format(DateFormat)
        } else {
            task.Date = next
        }
    }
    return nil
}