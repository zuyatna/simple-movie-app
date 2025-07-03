package repository

import (
	"database/sql"
	"movie-api/internal/model"

	"gorm.io/gorm"
)

type MovieRepository interface {
	FindAll() ([]model.Movie, error)
	FindByID(id uint) (model.Movie, error)
	Create(movie *model.Movie) (*model.Movie, error)
	Update(movie *model.Movie) (*model.Movie, error)
	Delete(id uint) error
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}

func (m movieRepository) FindAll() ([]model.Movie, error) {
	var movies []model.Movie
	err := m.db.Order("created_at DESC").Find(&movies).Error
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (m movieRepository) FindByID(id uint) (model.Movie, error) {
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

func (m movieRepository) Create(movie *model.Movie) (*model.Movie, error) {
	err := m.db.Create(movie).Error
	if err != nil {
		return nil, err
	}
	if movie.ID == 0 {
		return nil, sql.ErrNoRows
	}
	return movie, nil
}

func (m movieRepository) Update(movie *model.Movie) (*model.Movie, error) {
	err := m.db.Save(movie).Error
	if err != nil {
		return nil, err
	}
	if movie.ID == 0 {
		return nil, sql.ErrNoRows
	}
	return movie, nil
}

func (m movieRepository) Delete(id uint) error {
	return m.db.Delete(&model.Movie{}, id).Error
}
