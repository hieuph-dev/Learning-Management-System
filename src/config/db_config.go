package config

import (
	"fmt"
	"lms/src/utils"
)

type ServerConfig struct { // Bản thiết kết của một cái hộp
	ServerAddress string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type DBConfig struct {
	DB DatabaseConfig
}

// Máy sản xuất hộp theo bản thiết kế có sẵn
func NewServerConfig() *ServerConfig {
	return &ServerConfig{ // Trả về địa chỉ của cái hộp vừa tạo
		ServerAddress: ":8080",
	}
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		DB: DatabaseConfig{
			Host:     utils.GetEnv("DB_HOST", "postgres"),
			Port:     utils.GetEnv("DB_PORT", "5432"),
			User:     utils.GetEnv("DB_USER", "postgres"),
			Password: utils.GetEnv("DB_PASSWORD", "0917958087"),
			DBName:   utils.GetEnv("DB_NAME", "learning_management_system"),
			SSLMode:  utils.GetEnv("DB_SSLMODE", "disable"),
		},
	}
}

func (dbc *DBConfig) DNS() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbc.DB.Host, dbc.DB.Port, dbc.DB.User, dbc.DB.Password, dbc.DB.DBName, dbc.DB.SSLMode)
}
