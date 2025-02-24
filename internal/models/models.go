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
	Genre       string      `gorm:"type:varchar(100);not null"`
	ReleaseDate time.Time   `gorm:"not null"`
	PosterURL   string      `gorm:"type:varchar(255)"`
	ShowTimes   []ShowTime  `gorm:"foreignKey:MovieID"`
}



type Hall struct {
	gorm.Model
	Name      string     `gorm:"not null"`
	Capacity  int        `gorm:"not null"`
	ShowTimes []ShowTime `gorm:"foreignKey:HallID"`
}

type Seat struct {
	gorm.Model
	HallID     uint    `gorm:"not null;uniqueIndex:idx_hall_seat"`
	ShowTimeID uint    `gorm:"not null"`
	RowNumber  string  `gorm:"not null;type:varchar(2);uniqueIndex:idx_hall_seat"` // e.g., "A", "B"
	SeatNumber int     `gorm:"not null;uniqueIndex:idx_hall_seat"`
	Status     string  `gorm:"not null;type:varchar(20);default:'AVAILABLE'"` // AVAILABLE, RESERVED, BOOKED
	Tickets    []Ticket `gorm:"foreignKey:SeatID"`
}

type Ticket struct {
	gorm.Model
	UserID      uint      `gorm:"not null"`
	ShowTimeID  uint      `gorm:"not null"`
	SeatID      uint      `gorm:"not null;uniqueIndex:idx_showtime_seat"`
	Status      string    `gorm:"not null;type:varchar(20);default:'reserved'"` // reserved, paid, cancelled
	BookingCode string    `gorm:"not null;type:varchar(50);uniqueIndex"`
	Price       float64   `gorm:"not null;type:decimal(10,2)"`
}