package database

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host, Port, User, Pass, Name, Params string
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// LoadConfigFromEnv đọc cấu hình DB từ biến môi trường (.env)
func LoadConfigFromEnv() Config {
	_ = godotenv.Load() // không lỗi nếu thiếu file
	return Config{
		Host:   getEnv("DB_HOST", "127.0.0.1"),
		Port:   getEnv("DB_PORT", "3306"),
		User:   getEnv("DB_USER", "root"),
		Pass:   os.Getenv("DB_PASS"),
		Name:   getEnv("DB_NAME", "crud_api_user"),
		Params: getEnv("DB_PARAMS", "charset=utf8mb4&parseTime=True&loc=Local"),
	}
}

func Open(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name, cfg.Params)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}
