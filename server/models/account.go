package models

import (
  "time"
)

type Account struct {
    ID                 uint      `gorm:"primaryKey;autoIncrement"`
    UserID             uint      `gorm:"not null"`                 // Foreign key to User
    User               User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Relationship
    CardID             uint      `gorm:"not null"`                 // Foreign key to Card
    Card               Card      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Relationship
    AccountType        string    `gorm:"size:20;not null"`         // "Savings" or "Current"
    Balance            float64   `gorm:"type:decimal(15,2)"`       // Account balance
    DailyWithdrawalLimit float64 `gorm:"type:decimal(15,2)"`       // Daily withdrawal limit
    CreatedAt          time.Time `gorm:"autoCreateTime"`
    UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

