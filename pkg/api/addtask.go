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
        writeJson(w, map[string]string{"error": "ошибка десериализации JSON"})
        return
    }
    if task.Title == "" {
        writeJson(w, map[string]string{"error": "не указан заголовок задачи"})
        return
    }
    if err := checkDate(&task); err != nil {
        writeJson(w, map[string]string{"error": err.Error()})
        return
    }

    id, err := db.AddTask(&task)
    if err != nil {
        writeJson(w, map[string]string{"error": err.Error()})
        return
    }
    writeJson(w, map[string]interface{}{"id": id})
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

// // сравнение только дат (без учёта времени)
// func afterNow(date, now time.Time) bool {
//     y1, m1, d1 := date.Date()
//     y2, m2, d2 := now.Date()
//     t1 := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
//     t2 := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
//     return t1.After(t2)
// }