package models

import (
	"database/sql"
	"time"
)

type Email struct {
	Uuid         string          `json:"uuid" gorm:"primaryKey;column:uuid"`
	Email        string          `json:"email" gorm:"column:email"`
	Text         string          `json:"text" gorm:"column:text"`
	Provider     string          `json:"provider" gorm:"column:provider"`
	Status       string          `json:"status" gorm:"column:status"`
	ErrorDetails *sql.NullString `json:"error_details" gorm:"column:error_details"`
	CreatedAt    time.Time       `json:"created_at" gorm:"column:created_at"`
}

func (Email) TableName() string {
	return "email"
}

type ProviderEmail struct {
	Uuid    string `json:"uuid"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}
