package models

import (
	"time"

	"gorm.io/gorm"
)

type ShowTime struct {
	gorm.Model
	MovieID   uint      `gorm:"not null"`
	Movie     Movie     `gorm:"foreignKey:MovieID"`
	HallID    uint      `gorm:"not null"`
	Hall      Hall      `gorm:"foreignKey:HallID"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	Seats     []Seat    `gorm:"foreignKey:ShowTimeID"`
	Bookings  []Booking `gorm:"foreignKey:ShowTimeID"`
}