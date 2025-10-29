package server
import (
	"fmt"
	"net/http"
	"os"
	"planner/pkg/api"
)

func Run() {
    // Инициализация API обработчиков
    api.Init()

    // Определяем директорию фронтенда
    webDir := "./web"

    // Определяем порт
    port := os.Getenv("TODO_PORT")
    if port == "" {
        port = "7540"
    }

    // Создаём обработчик для отдачи статических файлов
    http.Handle("/", http.FileServer(http.Dir(webDir)))

    // Запускаем сервер
    addr := fmt.Sprintf(":%s", port)
    fmt.Printf("Сервер запущен на %s\n", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        panic(err)
    }
}