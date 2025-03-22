package models

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

// MusicInfo структура информации о песне.
// @Description Информация о песне, включая группу, название песни, дату создания, текст, ссылку.
type MusicInfo struct {
	gorm.Model
	Group       string `json:"group" gorm:"not null"`
	Song        string `json:"song" gorm:"not null"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text" gorm:"not null"`
	Link        string `json:"link"`
}

// Validate проверяет заполнение полей
func (m *MusicInfo) Validate() error {
	if strings.TrimSpace(m.Group) == "" {
		return errors.New("поле 'Group' обязательно для заполнения")
	}
	if strings.TrimSpace(m.Song) == "" {
		return errors.New("поле 'Song' обязательно для заполнения")
	}
	if strings.TrimSpace(m.Text) == "" {
		return errors.New("поле 'Text' обязательно для заполнения")
	}
	return nil
}
