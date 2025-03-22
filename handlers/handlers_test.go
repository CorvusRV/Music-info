package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"music-info/database"
	"music-info/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSongCreateHandler(t *testing.T) {

	// Создаем тестовые данные
	tx := database.SetupTestDB(t)
	defer database.DropTableDB(t, tx, "music_infos")

	songInfo := models.MusicInfo{
		Group:       "Muse",
		Song:        "Supermassive Black Hole",
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	body, _ := json.Marshal(songInfo)

	// Проверяем метод
	req, err := http.NewRequest("POST", "/songs/add", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/songs/add", SongCreateHandler).Methods("POST")
	router.ServeHTTP(rec, req)

	var response models.MusicInfo
	assert.Equal(t, http.StatusCreated, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, songInfo.Group, response.Group)
	assert.Equal(t, songInfo.Song, response.Song)
	assert.Equal(t, songInfo.ReleaseDate, response.ReleaseDate)
	assert.Equal(t, songInfo.Text, response.Text)
	assert.Equal(t, songInfo.Link, response.Link)
	assert.NotEqual(t, 0, response.ID)
}

func TestSongCreateHandlerBaseParams(t *testing.T) {

	// Создаем тестовые данные
	tx := database.SetupTestDB(t)
	defer database.DropTableDB(t, tx, "music_infos")

	songInfo := models.MusicInfo{
		Group: "Muse",
		Song:  "Supermassive Black Hole",
		Text:  "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
	}
	body, _ := json.Marshal(songInfo)

	// Проверяем метод
	req, err := http.NewRequest("POST", "/songs/add", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/songs/add", SongCreateHandler).Methods("POST")
	router.ServeHTTP(rec, req)

	var response models.MusicInfo
	assert.Equal(t, http.StatusCreated, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, songInfo.Group, response.Group)
	assert.Equal(t, songInfo.Song, response.Song)
	assert.Equal(t, songInfo.ReleaseDate, response.ReleaseDate)
	assert.Equal(t, songInfo.Text, response.Text)
	assert.Equal(t, songInfo.Link, response.Link)
	assert.NotEqual(t, 0, response.ID)
}

func TestSongDetailHandler(t *testing.T) {

	// Создаем тестовые данные
	tx := database.SetupTestDB(t)
	defer database.DropTableDB(t, tx, "music_infos")

	songInfo := models.MusicInfo{
		Group:       "Muse",
		Song:        "Supermassive Black Hole",
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	err := database.DBSongCreate(&songInfo)
	assert.NoError(t, err)

	// Проверяем метод
	req, err := http.NewRequest("GET", "/songs/info?group=Muse&song=Supermassive%20Black%20Hole", nil)
	assert.NoError(t, err)

	record := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/songs/info", SongDetailHandler).Methods("GET")
	router.ServeHTTP(record, req)

	var response models.MusicInfo
	assert.Equal(t, http.StatusOK, record.Code)
	err = json.Unmarshal(record.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, songInfo.ReleaseDate, response.ReleaseDate)
	assert.Equal(t, songInfo.Text, response.Text)
	assert.Equal(t, songInfo.Link, response.Link)
}

func TestSongDetailHandlerNotFound(t *testing.T) {

	// Создаем тестовые данные
	tx := database.SetupTestDB(t)
	defer database.DropTableDB(t, tx, "music_infos")

	songInfo := models.MusicInfo{
		Group:       "Muse",
		Song:        "Supermassive Black Hole",
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	err := database.DBSongCreate(&songInfo)
	assert.NoError(t, err)

	// Проверяем метод
	req, err := http.NewRequest("GET", "/songs/info?group=Muse&song=Supermassive%20Black%20Hole", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/songs/info", SongDetailHandler).Methods("GET")
	router.ServeHTTP(rec, req)

	var response map[string]string
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Error(t, err)
	assert.Equal(t, "", response["error"])
}

func TestSongUpdateHandler(t *testing.T) {

	// Создаем тестовые данные
	tx := database.SetupTestDB(t)
	defer database.DropTableDB(t, tx, "music_infos")

	songInfo := models.MusicInfo{
		Group:       "Muse",
		Song:        "Supermassive Black Hole",
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	err := database.DBSongCreate(&songInfo)
	assert.NoError(t, err)

	// Проверяем метод
	newText := "I thought I was a fool for no-one\nOh baby I'm a fool for you\nYou're the queen of the superficial\nAnd how long before you tell the truth\n\nOooh...You set my soul alight\nOooh...You set my soul alight"
	updateSong := models.MusicInfo{
		Text: newText,
	}
	body, _ := json.Marshal(updateSong)

	req, err := http.NewRequest("PUT", "/songs/info/update?group=Muse&song=Supermassive%20Black%20Hole", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/songs/info/update", SongUpdateHandler).Methods("PUT")
	router.ServeHTTP(rec, req)

	var response models.MusicInfo
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updateSong.Text, response.Text)

	result, err := database.DBSongDetail("Muse", "Supermassive Black Hole")
	assert.NoError(t, err)
	assert.Equal(t, updateSong.Text, result.Text)
}

func TestSongUpdateHandlerNotFound(t *testing.T) {

	// Создаем тестовые данные
	tx := database.SetupTestDB(t)
	defer database.DropTableDB(t, tx, "music_infos")

	// Проверяем метод
	newText := "I thought I was a fool for no-one\nOh baby I'm a fool for you\nYou're the queen of the superficial\nAnd how long before you tell the truth\n\nOooh...You set my soul alight\nOooh...You set my soul alight"
	updateSong := models.MusicInfo{
		Text: newText,
	}
	body, _ := json.Marshal(updateSong)

	req, err := http.NewRequest("PUT", "/songs/info/update?group=Queen&song=Bohemian Rhapsody", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/songs/info/update", SongUpdateHandler).Methods("PUT")
	router.ServeHTTP(rec, req)

	var response map[string]string
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Error(t, err)
	assert.Equal(t, "", response["error"])
}

func TestGetSongsHandler(t *testing.T) {

	// Создаем тестовые данные
	tx := database.SetupTestDB(t)
	defer database.DropTableDB(t, tx, "music_infos")

	songsInfo := []models.MusicInfo{
		{
			Group:       "Muse",
			Song:        "Supermassive Black Hole",
			ReleaseDate: "16.07.2006",
			Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
			Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		}, {
			Group:       "Queen",
			Song:        "Bohemian Rhapsody",
			ReleaseDate: "31.10.1975",
			Text:        "Is this the real life? Is this just fantasy?\nCaught in a landslide, no escape from reality\nOpen your eyes, look up to the skies and see\nI'm just a poor boy, I need no sympathy\nBecause I'm easy come, easy go\nLittle high, little low\nAny way the wind blows doesn't really matter to me, to me\n\nMama, just killed a man\nPut a gun against his head, pulled my trigger, now he's dead\nMama, life had just begun\nBut now I've gone and thrown it all away\nMama, ooh, didn't mean to make you cry\nIf I'm not back again this time tomorrow\nCarry on, carry on as if nothing really matters",
			Link:        "https://www.youtube.com/watch?v=vbvyNnw8Qjg",
		},
	}
	for _, msg := range songsInfo {
		err := database.DBSongCreate(&msg)
		assert.NoError(t, err)
	}

	// Проверяем метод
	req, err := http.NewRequest("GET", "/songs?group=Muse&page=1&limit=1", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/songs", GetSongsHandler).Methods("GET")
	router.ServeHTTP(rec, req)

	var response []models.MusicInfo
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, "Muse", response[0].Group)
}

func TestDeleteSongHandler(t *testing.T) {

	// Создаем тестовые данные
	tx := database.SetupTestDB(t)
	defer database.DropTableDB(t, tx, "music_infos")

	testData := models.MusicInfo{
		Group:       "Muse",
		Song:        "Supermassive Black Hole",
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	result := database.DB.Create(&testData)
	assert.NoError(t, result.Error)

	// Проверяем метод
	req, err := http.NewRequest("DELETE", "/songs/info/delete?group=Muse&song=Supermassive%20Black%20Hole", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/s", SongDeleteHandler).Methods("DELETE")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	var songDetail models.MusicInfo
	result = database.DB.Where("\"group\" = ? AND \"song\" = ?", "Muse", "Supermassive Black Hole").First(&songDetail)
	assert.Nil(t, result.Error)
	assert.NoError(t, result.Error)
}

func TestDeleteSongHandlerNotFound(t *testing.T) {

	// Создаем тестовые данные
	tx := database.SetupTestDB(t)
	defer database.DropTableDB(t, tx, "music_infos")

	// Проверяем метод
	req, err := http.NewRequest("DELETE", "/songs/info/delete?group=Muse&song=Supermassive Black Hole", nil)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/songs/info/delete", SongDeleteHandler).Methods("DELETE")
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "запись не найдена")
}
