package api

import (
    "errors"
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "time"
)

const DateFormat = "20060102"

// Обработчик GET /api/nextdate
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
    nowStr := r.FormValue("now")
    dateStr := r.FormValue("date")
    repeatStr := r.FormValue("repeat")

    if r.Method != http.MethodGet {
        http.Error(w, "Недопустимый метод", http.StatusMethodNotAllowed)
        return
    }

    var now time.Time
    var err error

    if nowStr == "" {
        now = time.Now()
    } else {
        now, err = time.Parse(DateFormat, nowStr)
        if err != nil {
            http.Error(w, "Некорректный формат параметра now", http.StatusBadRequest)
            return
        }
    }

    if dateStr == "" {
        http.Error(w, "Параметр date обязателен", http.StatusBadRequest)
        return
    }
    if repeatStr == "" {
        http.Error(w, "Параметр repeat обязателен", http.StatusBadRequest)
        return
    }

    nextDate, err := NextDate(now, dateStr, repeatStr)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Write([]byte(nextDate))
}

// Функция вычисления следующей даты
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
    if repeat == "" {
        return "", errors.New("пустое правило repeat")
    }

    date, err := time.Parse(DateFormat, dstart)
    if err != nil {
        return "", fmt.Errorf("ошибка парсинга даты: %v", err)
    }

    parts := strings.Split(strings.TrimSpace(repeat), " ")

    if len(parts) == 0 {
        return "", errors.New("некорректный формат repeat")
    }

    switch parts[0] {
    case "d":
        if len(parts) != 2 {
            return "", errors.New("не указан интервал для d")
        }

        interval, err := strconv.Atoi(parts[1])
        if err != nil || interval <= 0 || interval > 400 {
            return "", fmt.Errorf("некорректное значение интервала: %v", parts[1])
        }

        for {
            date = date.AddDate(0, 0, interval)
            if afterNow(date, now) {
                break
            }
        }
        return date.Format(DateFormat), nil

    case "y":
        for {
            date = date.AddDate(1, 0, 0)
            if afterNow(date, now) {
                break
            }
        }
        return date.Format(DateFormat), nil

    default:
        return "", fmt.Errorf("неподдерживаемое правило повторения: %v", repeat)
    }
}

// Вспомогательная функция сравнения только дат
func afterNow(date, now time.Time) bool {
    y1, m1, d1 := date.Date()
    y2, m2, d2 := now.Date()
    d1t := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
    d2t := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
    return d1t.After(d2t)
}
