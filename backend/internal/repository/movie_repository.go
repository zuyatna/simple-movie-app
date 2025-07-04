package repository

import (
	"database/sql"
	"movie-api/internal/model"

	"gorm.io/gorm"
)

type MovieRepository interface {
	FindAllMovie() ([]model.Movie, error)
	FindByMovieID(id uint) (model.Movie, error)
	CreateMovie(movie *model.Movie) (*model.Movie, error)
	UpdateMovie(movie *model.Movie) (*model.Movie, error)
	DeleteMovie(id uint) error
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}

func (m movieRepository) FindAllMovie() ([]model.Movie, error) {
	var movies []model.Movie
	err := m.db.Order("created_at DESC").Find(&movies).Error
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (m movieRepository) FindByMovieID(id uint) (model.Movie, error) {
	var movie model.Movie
	err := m.db.First(&movie, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.Movie{}, sql.ErrNoRows
		}
		return model.Movie{}, err
	}
	if movie.ID == 0 {
		return model.Movie{}, sql.ErrNoRows
	}
	return movie, nil
}

func (m movieRepository) CreateMovie(movie *model.Movie) (*model.Movie, error) {
	err := m.db.Create(movie).Error
	if err != nil {
		return nil, err
	}
	if movie.ID == 0 {
		return nil, sql.ErrNoRows
	}
	return movie, nil
}

func (m movieRepository) UpdateMovie(movie *model.Movie) (*model.Movie, error) {
	err := m.db.Save(movie).Error
	if err != nil {
		return nil, err
	}
	if movie.ID == 0 {
		return nil, sql.ErrNoRows
	}
	return movie, nil
}

func (m movieRepository) DeleteMovie(id uint) error {
	return m.db.Delete(&model.Movie{}, id).Error
}
