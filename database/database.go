package database

import (
	"errors"
	"fmt"
	"log"

	"music-info/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Инициализация базы данных
func InitDB(DNS string) {

	var err error
	DB, err = gorm.Open(postgres.Open(DNS), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка при открытии базы данных: %v", err)
	}

	err = DB.AutoMigrate(&models.MusicInfo{})
	if err != nil {
		log.Fatalf("Ошибка при создании таблицы: %v", err)
	}
	log.Println(DB.Name())
	log.Println("База данных инициализирована")

}

// Создание новой запись в базе данных.
func DBSongCreate(songInfo *models.MusicInfo) error {

	result := DB.Create(songInfo)

	return result.Error
}

// Обновление информации о песне по полям Group и Song.
func DBSongUpdate(group, song string, updateSong *models.MusicInfo) error {

	result := DB.Model(&models.MusicInfo{}).Where("\"group\" = ? AND \"song\" = ?", group, song).Updates(updateSong)

	return result.Error
}

// Удаление информации о песне по полям Group и Song.
func DBSongDelete(group, song string) error {

	result := DB.Where("\"group\" = ? AND \"song\" = ?", group, song).Delete(&models.MusicInfo{})

	if result.Error != nil {
		return fmt.Errorf("ошибка при удалении записи: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("запись не найдена: group=%s, song=%s", group, song)
	}

	return result.Error
}

// Возвращение информации о песне по полям Group и Song.
func DBSongDetail(group, song string) (*models.MusicInfo, error) {

	var songInfo models.MusicInfo

	if DB == nil {
		return nil, errors.New("база данных не инициализирована")
	}

	result := DB.Select("group", "song", "release_date", "text", "link").Where("\"group\" = ? AND \"song\" = ?", group, song).First(&songInfo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("запись не найдена: group=%s, song=%s", group, song)
		}
		return nil, result.Error
	}

	return &songInfo, nil
}

// Возвращение списка песен
func DBGetSongs(group string, page, limit int) ([]models.MusicInfo, error) {
	var songs []models.MusicInfo

	offset := (page - 1) * limit
	query := DB.Select("group", "song", "release_date", "text", "link").Where("\"group\" LIKE ?", "%"+group+"%").Offset(offset).Limit(limit).Order("\"group\", song, release_date, text, link").Find(&songs)

	return songs, query.Error
}
