package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"music-info/config"
	"music-info/database"
	"music-info/models"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func InitConfig(t *testing.T) config.Config {

	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	conf := config.Config{
		Port: "8080",
		DSN: "host=" + os.Getenv("DB_HOST") +
			" user=" + os.Getenv("DB_USER") +
			" dbname=testdb_test" +
			" password=" + os.Getenv("DB_PASSWORD") +
			" port=" + os.Getenv("DB_PORT") +
			" sslmode=" + os.Getenv("DB_SSLMODE") +
			" TimeZone=UTC",
	}

	return conf
}

func TestInitServer(t *testing.T) {

	config := InitConfig(t)

	// Проверяем инициализацию сервера
	server := initServer(config)

	assert.NotNil(t, server)

	// Проверяем, что адрес сервера корректный
	assert.Equal(t, ":8080", server.Addr)

	assert.NotNil(t, database.DB)
}

func TestServerEndpoints(t *testing.T) {

	// Инициализируем сервер
	config := InitConfig(t)
	server := initServer(config)

	songInfo := models.MusicInfo{
		Group:       "Muse",
		Song:        "Supermassive Black Hole",
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	result := database.DB.Create(&songInfo)
	assert.NoError(t, result.Error)

	// Создаем HTTP-клиент и проверяем ответ
	client := server.Handler

	req, err := http.NewRequest("GET", "/songs/info?group=Muse&song=Supermassive%20Black%20Hole", nil)
	assert.NoError(t, err)

	record := httptest.NewRecorder()
	client.ServeHTTP(record, req)

	assert.Equal(t, http.StatusOK, record.Code)

	var response models.MusicInfo
	err = json.Unmarshal(record.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Muse", response.Group)
	assert.Equal(t, "Supermassive Black Hole", response.Song)
}

func TestServerEndpointsNotFound(t *testing.T) {

	// Инициализируем сервер
	config := InitConfig(t)
	server := initServer(config)

	// Создаем HTTP-клиент и проверяем ответ
	client := server.Handler

	req, err := http.NewRequest("GET", "/endpoint", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	client.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
