package services

import (
	"context"
	"errors"
	"fmt"
	"movie-ticket-booking/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BookingService struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewBookingService(db *gorm.DB, redisClient *redis.Client) *BookingService {
	return &BookingService{
		db:          db,
		redisClient: redisClient,
	}
}

// CreateBooking creates a new booking with seat locking
func (s *BookingService) CreateBooking(ctx context.Context, userID uint, showtimeID uint, seatIDs []uint) (*models.Booking, error) {
	// Start a database transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get showtime details
	var showtime models.ShowTime
	if err := tx.First(&showtime, showtimeID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("showtime not found")
	}

	// Check if showtime has already started
	if time.Now().After(showtime.StartTime) {
		tx.Rollback()
		return nil, errors.New("cannot book seats for a show that has already started")
	}

	// Verify seats exist and are available
	var seats []*models.Seat
	for _, seatID := range seatIDs {
		// Check if seat is already locked
		lockKey := fmt.Sprintf("seat_lock:%d:%d", showtimeID, seatID)
		existingLock, err := s.redisClient.Get(ctx, lockKey).Result()
		if err != nil && err != redis.Nil {
			s.releaseSeatLocks(ctx, showtimeID, seatIDs)
			tx.Rollback()
			return nil, fmt.Errorf("error checking seat lock: %v", err)
		}
		if existingLock != "" {
			s.releaseSeatLocks(ctx, showtimeID, seatIDs)
			tx.Rollback()
			return nil, fmt.Errorf("seat %d is currently being booked by user %s", seatID, existingLock)
		}

		// Try to lock the seat in Redis
		locked, err := s.redisClient.SetNX(ctx, lockKey, userID, 5*time.Minute).Result()
		if err != nil || !locked {
			// Release any locks we've acquired
			s.releaseSeatLocks(ctx, showtimeID, seatIDs)
			tx.Rollback()
			return nil, fmt.Errorf("failed to lock seat %d: %v", seatID, err)
		}

		var seat models.Seat
		if err := tx.First(&seat, seatID).Error; err != nil {
			s.releaseSeatLocks(ctx, showtimeID, seatIDs)
			tx.Rollback()
			return nil, fmt.Errorf("seat %d not found: %v", seatID, err)
		}

		// Verify seat belongs to the correct showtime and is available
		if seat.ShowTimeID != showtimeID {
			s.releaseSeatLocks(ctx, showtimeID, seatIDs)
			tx.Rollback()
			return nil, fmt.Errorf("seat %d belongs to showtime %d, not %d", seatID, seat.ShowTimeID, showtimeID)
		}
		if seat.Status != models.SeatStatusAvailable {
			s.releaseSeatLocks(ctx, showtimeID, seatIDs)
			tx.Rollback()
			return nil, fmt.Errorf("seat %d is not available (current status: %s)", seatID, seat.Status)
		}

		seats = append(seats, &seat)
	}

	// Calculate total amount
	totalAmount := float64(len(seatIDs)) * showtime.Price

	// Create booking
	booking := &models.Booking{
		UserID:      userID,
		ShowTimeID:  showtimeID,
		TotalAmount: totalAmount,
		Status:      models.BookingStatusConfirmed,
		BookedAt:    time.Now(),
	}

	if err := tx.Create(booking).Error; err != nil {
		s.releaseSeatLocks(ctx, showtimeID, seatIDs)
		tx.Rollback()
		return nil, err
	}

	// Update seat status and create booking seats
	for _, seat := range seats {
		seat.Status = models.SeatStatusBooked
		if err := tx.Save(seat).Error; err != nil {
			s.releaseSeatLocks(ctx, showtimeID, seatIDs)
			tx.Rollback()
			return nil, err
		}

		// Create booking seats relationship
		bookingSeat := &models.BookingSeat{
			BookingID: booking.ID,
			SeatID:    seat.ID,
			Price:     showtime.Price,
		}
		if err := tx.Create(bookingSeat).Error; err != nil {
			s.releaseSeatLocks(ctx, showtimeID, seatIDs)
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		s.releaseSeatLocks(ctx, showtimeID, seatIDs)
		return nil, err
	}

	// Release locks after successful commit
	s.releaseSeatLocks(ctx, showtimeID, seatIDs)

	return booking, nil
}

// GetBooking retrieves a booking by ID
func (s *BookingService) GetBooking(id uint) (*models.Booking, error) {
	var booking models.Booking
	if err := s.db.Preload("User").Preload("Showtime").Preload("Seats").First(&booking, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}
	return &booking, nil
}

// GetUserBookings retrieves all bookings for a user
func (s *BookingService) GetUserBookings(userID uint) ([]*models.Booking, error) {
	var bookings []*models.Booking
	if err := s.db.Where("user_id = ?", userID).Preload("Showtime").Preload("Seats").Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

// CancelBooking cancels a booking and releases the seats
func (s *BookingService) CancelBooking(ctx context.Context, bookingID uint, userID uint) error {
	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get booking with seats
	var booking models.Booking
	if err := tx.Preload("Seats.Seat").First(&booking, bookingID).Error; err != nil {
		tx.Rollback()
		return errors.New("booking not found")
	}

	// Verify booking belongs to user
	if booking.UserID != userID {
		tx.Rollback()
		return errors.New("unauthorized to cancel this booking")
	}

	// Check if booking can be cancelled (e.g., not already cancelled)
	if booking.Status == models.BookingStatusCancelled {
		tx.Rollback()
		return errors.New("booking is already cancelled")
	}

	// Update booking status
	booking.Status = models.BookingStatusCancelled
	if err := tx.Save(&booking).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Release seats
	for _, bookingSeat := range booking.Seats {
		bookingSeat.Seat.Status = models.SeatStatusAvailable
		if err := tx.Save(&bookingSeat.Seat).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	return tx.Commit().Error
}

// Helper function to release seat locks in Redis
func (s *BookingService) releaseSeatLocks(ctx context.Context, showtimeID uint, seatIDs []uint) {
	for _, seatID := range seatIDs {
		lockKey := fmt.Sprintf("seat_lock:%d:%d", showtimeID, seatID)
		s.redisClient.Del(ctx, lockKey)
	}
}