package database

import (
	"log"
	"testing"

	"music-info/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDBSongCreate(t *testing.T) {

	// Создаем тестовые данные
	tx := SetupTestDB(t)
	defer DropTableDB(t, tx, "music_infos")

	songInfo := models.MusicInfo{
		Group:       "Muse",
		Song:        "Supermassive Black Hole",
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	// Проверяем метод
	err := DBSongCreate(&songInfo)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, songInfo.ID)
}

func TestDBSongUpdate(t *testing.T) {

	// Создаем тестовые данные
	tx := SetupTestDB(t)
	defer DropTableDB(t, tx, "music_infos")

	songInfo := models.MusicInfo{
		Group:       "Muse",
		Song:        "Supermassive Black Hole",
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	newText := "I thought I was a fool for no-one\nOh baby I'm a fool for you\nYou're the queen of the superficial\nAnd how long before you tell the truth\n\nOooh...You set my soul alight\nOooh...You set my soul alight"
	updateSong := models.MusicInfo{
		Text: newText,
	}

	// Проверяем метод
	err := DBSongCreate(&songInfo)
	assert.NoError(t, err)

	err = DBSongUpdate("Muse", "Supermassive Black Hole", &updateSong)
	assert.NoError(t, err)

	result, err := DBSongDetail("Muse", "Supermassive Black Hole")
	if err != nil {
		t.Fatalf("Ошибка при вызове DBSongDetail: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, newText, result.Text)
}

func TestDBDeleteSong(t *testing.T) {

	// Создаем тестовые данные
	tx := SetupTestDB(t)
	defer DropTableDB(t, tx, "music_infos")

	songInfo := models.MusicInfo{
		Group:       "Muse",
		Song:        "Supermassive Black Hole",
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	err := DBSongCreate(&songInfo)
	assert.NoError(t, err)

	// Проверяем метод
	err = DBSongDelete("Muse", "Supermassive Black Hole")
	assert.NoError(t, err)

	var songDetail models.MusicInfo
	result := DB.Where("\"group\" = ? AND \"song\" = ?", "Muse", "Supermassive Black Hole").First(&songDetail)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestDBDeleteSongNotFound(t *testing.T) {

	// Создаем тестовые данные
	tx := SetupTestDB(t)
	defer DropTableDB(t, tx, "music_infos")

	// Проверяем метод
	err := DBSongDelete("Muse", "Supermassive Black Hole")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "запись не найдена")
}

func TestDBSongDetail(t *testing.T) {

	// Создаем тестовые данные
	tx := SetupTestDB(t)
	defer DropTableDB(t, tx, "music_infos")

	songInfo := models.MusicInfo{
		Group:       "Muse",
		Song:        "Supermassive Black Hole",
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	err := DBSongCreate(&songInfo)
	assert.NoError(t, err)

	// Проверяем метод
	result, err := DBSongDetail("Muse", "Supermassive Black Hole")
	if err != nil {
		t.Fatalf("Ошибка при вызове метода DBSongDetail: %v", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, songInfo.Group, result.Group)
	assert.Equal(t, songInfo.Song, result.Song)
	assert.Equal(t, songInfo.ReleaseDate, result.ReleaseDate)
	assert.Equal(t, songInfo.Text, result.Text)
	assert.Equal(t, songInfo.Link, result.Link)
}

func TestDBGetSongs(t *testing.T) {

	// Создаем тестовые данные
	tx := SetupTestDB(t)
	defer DropTableDB(t, tx, "music_infos")

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
	for _, song := range songsInfo {
		err := DBSongCreate(&song)
		assert.NoError(t, err)
	}

	// Проверяем метод
	result, err := DBGetSongs("Muse", 1, 1)
	if err != nil {
		t.Fatalf("Ошибка при вызове метода DBGetSongs: %v", err)
	}
	log.Println(result)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "Muse", result[0].Group)
}
