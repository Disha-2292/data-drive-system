package models

import (
	"time"
)

type BaseModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// {
// 	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ4MjQyNDksInVzZXJfaWQiOjF9._fK8atPNZIe8uitqJGYleA498TBxlWYVx1TqaVAJeiw"
//   }
