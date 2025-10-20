package db

import (
	"fmt"
	"database/sql"
)

type Task struct {
    ID      string `json:"id,omitempty"`
    Date    string `json:"date"`
    Title   string `json:"title"`
    Comment string `json:"comment"`
    Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
    query := `INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?)`
    res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
    if err != nil {
        return 0, err
    }
    return res.LastInsertId()
}


func Tasks(limit int) ([]*Task, error) {
    query := `
        SELECT id, date, title, comment, repeat
        FROM scheduler
        ORDER BY date ASC
        LIMIT ?
    `

    rows, err := DB.Query(query, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var tasks []*Task

    for rows.Next() {
        var t Task
        var id int64

        if err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
            return nil, err
        }

        t.ID = fmt.Sprintf("%d", id)
        tasks = append(tasks, &t)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    if tasks == nil {
        tasks = []*Task{} // предотвращает {"tasks":null}
    }

    return tasks, nil
}

func GetTask(id string) (*Task, error) {
    query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
    var t Task
    var idInt int64
    err := DB.QueryRow(query, id).Scan(&idInt, &t.Date, &t.Title, &t.Comment, &t.Repeat)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("задача не найдена")
        }
        return nil, err
    }
    t.ID = fmt.Sprintf("%d", idInt)
    return &t, nil
}

func UpdateTask(task *Task) error {
    query := `
        UPDATE scheduler
        SET date = ?, title = ?, comment = ?, repeat = ?
        WHERE id = ?
    `
    res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
    if err != nil {
        return err
    }

    count, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if count == 0 {
        return fmt.Errorf("incorrect id for updating task")
    }
    return nil
}

// Удаляет задачу по идентификатору
func DeleteTask(id string) error {
    query := `DELETE FROM scheduler WHERE id = ?`
    res, err := DB.Exec(query, id)
    if err != nil {
        return err
    }
    n, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if n == 0 {
        return fmt.Errorf("задача не найдена")
    }
    return nil
}

// Обновляет только дату следующего выполнения
func UpdateDate(next, id string) error {
    query := `UPDATE scheduler SET date = ? WHERE id = ?`
    res, err := DB.Exec(query, next, id)
    if err != nil {
        return err
    }
    n, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if n == 0 {
        return fmt.Errorf("задача не найдена")
    }
    return nil
}