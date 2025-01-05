package models

import (
	"time"
)

type Transaction struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"`       // Maps to SERIAL PRIMARY KEY
	AccountID       uint      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Foreign key to accounts table
	Account         Account   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`          // Relationship with Account model
	TransactionType string    `gorm:"size:20;not null"`               // "withdrawal", "deposit", "transfer", "inquiry"
	Amount          float64   `gorm:"type:decimal(15,2)"`            // Maps to DECIMAL(15,2)
	Status          string    `gorm:"size:20;not null"`               // "success", "failed", "pending"
	ReferenceNumber string    `gorm:"size:50;unique;not null"`       // Maps to VARCHAR(50) UNIQUE
	Notes           string    `gorm:"type:text"`                     // Maps to TEXT
	CreatedAt       time.Time `gorm:"autoCreateTime"`                 // Maps to TIMESTAMP DEFAULT CURRENT_TIMESTAMP
}


type DailyLimit struct {
	ID               uint      `gorm:"primaryKey;autoIncrement"`        // Maps to SERIAL PRIMARY KEY
	AccountID        uint      `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Foreign key to accounts table
	Account          Account   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`          // Relationship with Account model
	Date             time.Time `gorm:"type:date;default:CURRENT_DATE"`   // Maps to DATE DEFAULT CURRENT_DATE
	TotalWithdrawal  float64   `gorm:"type:decimal(15,2);default:0.00"`  // Maps to DECIMAL(15,2) DEFAULT 0.00
	TotalTransactions int      `gorm:"default:0"`                        // Maps to INT DEFAULT 0
	CreatedAt        time.Time `gorm:"autoCreateTime"`                   // Maps to TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`                   // Maps to TIMESTAMP DEFAULT CURRENT_TIMESTAMP
}
