package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"movie-api/internal/model"
	"movie-api/internal/repository"
	"time"

	"github.com/go-redis/redis"
)

type MovieUsecase interface {
	FindAllMovie() ([]model.Movie, error)
	FindByMovieID(id uint) (model.Movie, error)
	CreateMovie(movie *model.Movie) (*model.Movie, error)
	UpdateMovie(id uint, movie *model.Movie) (*model.Movie, error)
	DeleteMovie(id uint) error
}

type movieUsecase struct {
	movieRepo   repository.MovieRepository
	redisClient *redis.Client
	ctx         context.Context
}

func NewMovieUsecase(movieRepo repository.MovieRepository, redisClient *redis.Client, ctx context.Context) MovieUsecase {
	return &movieUsecase{
		movieRepo:   movieRepo,
		redisClient: redisClient,
		ctx:         ctx,
	}
}

func (m *movieUsecase) FindAllMovie() ([]model.Movie, error) {
	cacheKey := "all_movies"
	val, err := m.redisClient.Get(cacheKey).Result()
	if err == redis.Nil {
		movies, err := m.movieRepo.FindAllMovie()
		if err != nil {
			return nil, err
		}

		cacheableMovies := make([]model.Movie, len(movies))
		for i, movie := range movies {
			cacheableMovies[i] = movie.ToCacheableMovie()
		}

		jsonData, _ := json.Marshal(cacheableMovies)
		m.redisClient.Set(cacheKey, jsonData, time.Hour*1).Err()
		return movies, nil
	}
	if err != nil {
		return nil, err
	}

	var movies []model.Movie
	if json.Unmarshal([]byte(val), &movies) == nil {
		return movies, nil
	}

	movies, err = m.movieRepo.FindAllMovie()
	if err != nil {
		return nil, err
	}

	cacheableMovies := make([]model.Movie, len(movies))
	for i, movie := range movies {
		cacheableMovies[i] = movie.ToCacheableMovie()
	}

	jsonData, _ := json.Marshal(cacheableMovies)
	m.redisClient.Set(cacheKey, jsonData, time.Hour*1).Err()
	return movies, nil
}

func (m *movieUsecase) FindByMovieID(id uint) (model.Movie, error) {
	cacheKey := fmt.Sprintf("movie:%d", id)
	val, err := m.redisClient.Get(cacheKey).Result()
	if err == redis.Nil {
		movie, err := m.movieRepo.FindByMovieID(id)
		if err != nil {
			return model.Movie{}, err
		}

		jsonData, _ := json.Marshal(movie)
		m.redisClient.Set(cacheKey, jsonData, time.Hour*1).Err()
		return movie, nil
	}
	if err != nil {
		return model.Movie{}, err
	}

	var movie model.Movie
	if json.Unmarshal([]byte(val), &movie) == nil {
		return movie, nil
	}

	movie, err = m.movieRepo.FindByMovieID(id)
	if err != nil {
		return model.Movie{}, err
	}

	jsonData, _ := json.Marshal(movie)
	m.redisClient.Set(cacheKey, jsonData, time.Hour*1).Err()
	return movie, nil
}

func (m *movieUsecase) CreateMovie(movie *model.Movie) (*model.Movie, error) {
	newMovie, err := m.movieRepo.CreateMovie(movie)
	if err != nil {
		return nil, err
	}
	m.redisClient.Del("all_movies").Err()
	return newMovie, nil
}

func (m *movieUsecase) UpdateMovie(id uint, movie *model.Movie) (*model.Movie, error) {
	movieData, err := m.movieRepo.FindByMovieID(id)
	if err != nil {
		return nil, err
	}

	movieData.Title = movie.Title
	movieData.Description = movie.Description
	movieData.Year = movie.Year
	movieData.Rating = movie.Rating
	if movie.Poster != nil {
		movieData.Poster = movie.Poster
	}

	updatedMovie, err := m.movieRepo.UpdateMovie(id, &movieData)
	if err != nil {
		return nil, err
	}

	m.redisClient.Del("all_movies").Err()
	m.redisClient.Del(fmt.Sprintf("movie:%d", id)).Err()
	return updatedMovie, nil
}

func (m *movieUsecase) DeleteMovie(id uint) error {
	err := m.movieRepo.DeleteMovie(id)
	if err == nil {
		m.redisClient.Del("all_movies").Err()
		m.redisClient.Del(fmt.Sprintf("movie:%d", id)).Err()
	}
	return err
}
