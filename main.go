package main

import (
	"log"
	"planner/pkg/db"
	"planner/pkg/server"
	_ "github.com/joho/godotenv/autoload"

)

func main() {
	dbFile := "scheduler.db"
// Инициализация базы данных
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
// Запуск сервера
    server.Run()

}