package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Guid         string    `gorm:"unique"`
	Username     string    `gorm:"type:varchar(15);unique;not null"`
	Email        string    `gorm:"type:varchar(20);unique;not null"`
	PasswordHash string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Role         string    `gorm:"type:varchar(10);default:'user'"`
}
