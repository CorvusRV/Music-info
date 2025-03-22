package database

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"

	"gorm.io/gorm"
)

// Инициализация тестовой БД
func initTestDB(t *testing.T) {

	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	DNSTest := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" dbname=testdb_test" +
		" password=" + os.Getenv("DB_PASSWORD") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=" + os.Getenv("DB_SSLMODE") +
		" TimeZone=UTC"

	InitDB(DNSTest)
}

// Настройка тестовой БД
func SetupTestDB(t *testing.T) *gorm.DB {

	initTestDB(t)

	tx := DB.Begin()
	if tx.Error != nil {
		t.Fatalf("Ошибка при начале транзакции: %v", tx.Error)
	}

	return tx
}

// Удаление таблицы в тестовой БД
func DropTableDB(t *testing.T, db *gorm.DB, table string) {

	err := DB.Migrator().DropTable(table)
	if err != nil {
		log.Printf("Ошибка удаления таблицы %s: %v\n", table, err)
	} else {
		log.Printf("Таблица %s успешно удалена\n", table)
	}
}
