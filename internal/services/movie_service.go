package services

import (
	"errors"
	"movie-ticket-booking/internal/models"
	"gorm.io/gorm"
)

type MovieService struct {
	db *gorm.DB
}

func NewMovieService(db *gorm.DB) *MovieService {
	return &MovieService{
		db: db,
	}
}

// GetMovies returns paginated movies from the database
func (s *MovieService) GetMovies(offset, limit int) ([]*models.Movie, int64, error) {
	var movies []*models.Movie
	var total int64

	// Get total count of movies
	if err := s.db.Model(&models.Movie{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated movies
	if err := s.db.Offset(offset).Limit(limit).Find(&movies).Error; err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

// GetMovieByID returns a specific movie by ID
func (s *MovieService) GetMovieByID(id uint) (*models.Movie, error) {
	var movie models.Movie
	if err := s.db.First(&movie, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("movie not found")
		}
		return nil, err
	}
	return &movie, nil
}

// CreateMovie creates a new movie in the database
func (s *MovieService) CreateMovie(movie *models.Movie) error {
	return s.db.Create(movie).Error
}

// UpdateMovie updates an existing movie
func (s *MovieService) UpdateMovie(movie *models.Movie) error {
	if err := s.db.First(&models.Movie{}, movie.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("movie not found")
		}
		return err
	}
	return s.db.Save(movie).Error
}

// DeleteMovie deletes a movie by ID
func (s *MovieService) DeleteMovie(id uint) error {
	result := s.db.Delete(&models.Movie{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("movie not found")
	}
	return nil
}