package server
import (
	"fmt"
	"net/http"
	"os"
	"planner/pkg/api"
)

// type MyServer struct {
// 	Logger *log.Logger
// 	Server http.Server
// }

// func CreateServer(log *log.Logger) *MyServer {
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", handlers.HandleMain)
// 	mux.HandleFunc("/upload", handlers.HandleUpload)
// 	return &MyServer{
// 		Server: http.Server{
// 			Addr:	":7540",
// 			Handler: mux,
// 			ErrorLog: log,
// 			ReadTimeout:  5 * time.Second,
// 			WriteTimeout: 10 * time.Second,
// 			IdleTimeout:  15 * time.Second,
// 		},
// 		Logger: log,
// 	}
// }

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