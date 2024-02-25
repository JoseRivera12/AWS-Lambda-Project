package store

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Email string `gorm:"unique;not null"`
}

type Transaction struct {
	ID          uint            `gorm:"primaryKey"`
	UserID      uint            `gorm:"not null"`
	User        User            `gorm:"foreignKey:UserID"`
	Amount      decimal.Decimal `gorm:"type:numeric(20,2)"`
	Description string
	Currency    string
	CreatedAt   time.Time `gorm:"not null"`
}
