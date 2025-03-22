// @title Music Info API
// @version 1.0
// @description API для управления информацией о музыке
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
package main

import (
	"log"
	"net/http"

	"music-info/config"
	"music-info/database"
	"music-info/handlers"

	_ "music-info/docs"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// initServer инициализирует и возвращает HTTP-сервер.
func initServer(config config.Config) *http.Server {
	// Инициализация логгера
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Загружаем переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	// Инициализация базы данных
	database.InitDB(config.DSN)

	// Настройка маршрутизатора
	router := mux.NewRouter()

	// Регистрируем обработчики
	router.HandleFunc("/songs/add", handlers.SongCreateHandler).Methods("POST")
	router.HandleFunc("/songs", handlers.GetSongsHandler).Methods("GET")
	router.HandleFunc("/songs/info", handlers.SongDetailHandler).Methods("GET")
	router.HandleFunc("/songs/info/update", handlers.SongUpdateHandler).Methods("PUT")
	router.HandleFunc("/songs/info/delete", handlers.SongDeleteHandler).Methods("DELETE")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Создаем HTTP-сервер
	return &http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}
}

func main() {
	config := config.LoadConfig()
	server := initServer(config)
	log.Printf("Сервер запущен на http://localhost%s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}
