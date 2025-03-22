package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageValidation(t *testing.T) {
	TestInfo := MusicInfo{
		Group:       "Muse",
		Song:        "Supermassive Black Hole",
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	// Проверка на пустые поля
	musicInfoGroup := TestInfo
	musicInfoGroup.Group = ""
	err := musicInfoGroup.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "поле 'Group' обязательно для заполнения")

	musicInfoSong := TestInfo
	musicInfoSong.Song = ""
	err = musicInfoSong.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "поле 'Song' обязательно для заполнения")

	musicInfoText := TestInfo
	musicInfoText.Text = ""
	err = musicInfoText.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "поле 'Text' обязательно для заполнения")

	// Проверка на корректные данные
	musicInfo := TestInfo
	err = musicInfo.Validate()
	assert.NoError(t, err)
}
