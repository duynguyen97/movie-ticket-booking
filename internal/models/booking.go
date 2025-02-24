package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID      uint      `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserID"`
	ShowTimeID  uint      `gorm:"not null"`
	Showtime    ShowTime  `gorm:"foreignKey:ShowTimeID"` // Fixed field name to match GraphQL schema
	TotalAmount float64   `gorm:"not null"`
	Status      string    `gorm:"not null"` // CONFIRMED, CANCELLED
	BookedAt    time.Time `gorm:"not null"`
	Seats       []BookingSeat
}

type BookingSeat struct {
	gorm.Model
	BookingID uint    `gorm:"not null"`
	SeatID    uint    `gorm:"not null"`
	Seat      Seat    `gorm:"foreignKey:SeatID"`
	Price     float64 `gorm:"not null"`
}

const (
	BookingStatusConfirmed = "CONFIRMED"
	BookingStatusCancelled = "CANCELLED"
	SeatStatusAvailable    = "AVAILABLE"
	SeatStatusReserved     = "RESERVED"
	SeatStatusBooked       = "BOOKED"
)