package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Name      string    `gorm:"not null"`
	Phone     string    `gorm:"not null"`
	Tickets   []Ticket  `gorm:"foreignKey:UserID"`
}

type Movie struct {
	gorm.Model
	Title       string      `gorm:"not null"`
	Description string      `gorm:"type:text"`
	Duration    int         `gorm:"not null"` // in minutes
	Genre       string      `gorm:"not null"`
	ShowTimes   []ShowTime  `gorm:"foreignKey:MovieID"`
}

type ShowTime struct {
	gorm.Model
	MovieID   uint      `gorm:"not null"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	HallID    uint      `gorm:"not null"`
	Seats     []Seat    `gorm:"foreignKey:ShowTimeID"`
}

type Hall struct {
	gorm.Model
	Name      string     `gorm:"not null"`
	Capacity  int        `gorm:"not null"`
	ShowTimes []ShowTime `gorm:"foreignKey:HallID"`
}

type Seat struct {
	gorm.Model
	ShowTimeID uint    `gorm:"not null"`
	Number     string  `gorm:"not null"` // e.g., "A1", "B2"
	Status     string  `gorm:"not null"` // available, booked, reserved
	Tickets    []Ticket `gorm:"foreignKey:SeatID"`
}

type Ticket struct {
	gorm.Model
	UserID     uint      `gorm:"not null"`
	SeatID     uint      `gorm:"not null"`
	Price      float64   `gorm:"not null"`
	BookedAt   time.Time `gorm:"not null"`
	Status     string    `gorm:"not null"` // booked, paid, cancelled
}