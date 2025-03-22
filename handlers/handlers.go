package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"music-info/database"
	"music-info/models"

	"gorm.io/gorm"
)

// sendError отправляет ошибку клиенту.
func sendError(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{Error: errorMessage})
}

// SongCreateHandler создает новое сообщение.
// @Summary Создать новое сообщение
// @Description Добавляет новое сообщение в базу данных
// @Tags messages
// @Accept json
// @Produce json
// @Param message body models.MusicInfo true "Данные сообщения"
// @Success 201 {object} models.MusicInfo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /messages [post]
func SongCreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var songInfo models.MusicInfo
	err := json.NewDecoder(r.Body).Decode(&songInfo)
	if err != nil {
		log.Printf("Ошибка при декодировании JSON: %v", err)
		sendError(w, http.StatusBadRequest, "Неверный формат JSON")
		return
	}

	if err := songInfo.Validate(); err != nil {
		log.Printf("Ошибка валидации: %v", err)
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = database.DBSongCreate(&songInfo)
	if err != nil {
		log.Printf("Ошибка при вставке данных: %v", err)
		sendError(w, http.StatusInternalServerError, "Ошибка при вставке данных")
		return
	}

	log.Printf("Новое сообщение добавлено в базу данных: %+v\n", songInfo)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(songInfo)
}

// SongDetailHandler возвращает информацию о песне.
// @Summary Получить информацию о песне
// @Description Возвращает информацию о песне по указанным группе и названию песни
// @Tags songs
// @Produce json
// @Param group query string true "Название группы"
// @Param song query string true "Название песни"
// @Success 200 {object} models.MusicInfo "Успешный ответ с информацией о песне"
// @Failure 400 {object} map[string]string "Неверные параметры запроса"
// @Failure 404 {object} map[string]string "Запись не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /songs/detail [get]
func SongDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")

	songInfo, err := database.DBSongDetail(group, song)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("Сообщение не найдено: group=%s, song=%s\n", group, song)
			sendError(w, http.StatusNotFound, "Сообщение не найдено")
		} else {
			log.Printf("Ошибка при получении сообщения: %v\n", err)
			sendError(w, http.StatusInternalServerError, "Ошибка при получении сообщения")
		}
		return
	}

	log.Printf("Сообщение найдено: %+v\n", songInfo)
	json.NewEncoder(w).Encode(songInfo)
}

// SongUpdateHandler обновляет информацию о песне.
// @Summary Обновить информацию о песне
// @Description Обновляет информацию о песне по указанным группе и названию песни
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string true "Название группы"
// @Param song query string true "Название песни"
// @Param updateInfo body models.MusicInfo true "Данные для обновления"
// @Success 200 {object} models.MusicInfo "Успешный ответ с обновлённой информацией о песне"
// @Failure 400 {object} map[string]string "Неверный формат JSON или параметры запроса"
// @Failure 404 {object} map[string]string "Запись не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /songs/update [put]
func SongUpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")

	var updateInfo models.MusicInfo
	err := json.NewDecoder(r.Body).Decode(&updateInfo)
	if err != nil {
		log.Printf("Ошибка при декодировании JSON: %v\n", err)
		sendError(w, http.StatusBadRequest, "Неверный формат JSON")
		return
	}

	err = database.DBSongUpdate(group, song, &updateInfo)
	if err != nil {
		log.Printf("Ошибка при обновлении сообщения: %v\n", err)
		sendError(w, http.StatusInternalServerError, "Ошибка при обновлении сообщения")
		return
	}

	log.Printf("Сообщение обновлено: group=%s, song=%s\n", group, song)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updateInfo)
}

// GetSongsHandler возвращает список песен.
// @Summary Получить список песен
// @Description Возвращает список песен с возможностью фильтрации по группе и пагинацией
// @Tags songs
// @Produce json
// @Param group query string false "Фильтр по группе"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество записей на странице" default(10)
// @Success 200 {array} models.MusicInfo "Успешный ответ со списком песен"
// @Failure 400 {object} map[string]string "Неверные параметры запроса"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /songs [get]
func GetSongsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	group := r.URL.Query().Get("group")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}

	messages, err := database.DBGetSongs(group, page, limit)
	if err != nil {
		log.Printf("Ошибка при получении данных: %v", err)
		sendError(w, http.StatusInternalServerError, "Ошибка при получении данных")
		return
	}

	log.Printf("Получено %d сообщений\n", len(messages))
	json.NewEncoder(w).Encode(messages)
}

// SongDeleteHandler удаляет запись о песне.
// @Summary Удалить запись о песне
// @Description Удаляет запись о песне по указанным группе и названию песни
// @Tags songs
// @Produce json
// @Param group query string true "Название группы"
// @Param song query string true "Название песни"
// @Success 204 "Запись успешно удалена"
// @Failure 400 {object} map[string]string "Неверные параметры запроса"
// @Failure 404 {object} map[string]string "Запись не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /songs [delete]
func SongDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")

	err := database.DBSongDelete(group, song)
	if err != nil {
		log.Printf("Ошибка при удалении записи: %v\n", err)
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Запись удалена: group=%s, song=%s\n", group, song)
	w.WriteHeader(http.StatusNoContent)
}
