package models

import (
	"time"
)
type Card struct {
    ID            uint      `gorm:"primaryKey;autoIncrement"`
    UserID        uint      `gorm:"not null"`
    User          User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    CardNumber    string    `gorm:"size:16;unique;not null"`
    PinHash       string    `gorm:"not null"`
    CardStatus    string    `gorm:"size:20;default:'active'"`
    FailedAttempts int      `gorm:"default:0"`
    ExpiryDate    time.Time `gorm:"not null"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

