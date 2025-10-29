package db

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // ← новый драйвер, вместо modernc.org/sqlite
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

var DB *sql.DB // Дополнительно для будущих операций с БД

func Init(dbFile string) error {
	var install bool
	if _, err := os.Stat(dbFile); errors.Is(err, os.ErrNotExist) {
		install = true
	}

	db, err := sql.Open("sqlite3", dbFile) // ← "sqlite3" для mattn/go-sqlite3
	if err != nil {
		return err
	}
	DB = db // Для удобства — глобальная переменная

	if install {
		if _, err := db.Exec(schema); err != nil {
			return err
		}
		log.Printf("Создана база и таблица scheduler")
	} else {
		log.Printf("Файл базы данных найден, создание таблицы не требуется")
	}
	return nil
}
