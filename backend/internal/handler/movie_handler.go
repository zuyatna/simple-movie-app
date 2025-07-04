package handler

import (
	"fmt"
	"io"
	"movie-api/internal/model"
	"movie-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	movieUsecase usecase.MovieUsecase
}

func NewMovieHandler(movieUsecase usecase.MovieUsecase) *MovieHandler {
	return &MovieHandler{
		movieUsecase: movieUsecase,
	}
}

func (mh *MovieHandler) FindAllMovie(c *gin.Context) {
	movies, err := mh.movieUsecase.FindAllMovie()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}

	for i := range movies {
		movies[i].PosterURL = fmt.Sprintf("/api/v1/movies/%d/poster", movies[i].ID)
	}

	c.JSON(http.StatusOK, movies)
}

func (mh *MovieHandler) FindByMovieID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	movie, err := mh.movieUsecase.FindByMovieID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	movie.PosterURL = fmt.Sprintf("/api/v1/movies/%d/poster", movie.ID)
	c.JSON(http.StatusOK, movie)
}

func (mh *MovieHandler) GetMoviePoster(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	movie, err := mh.movieUsecase.FindByMovieID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}
	if len(movie.Poster) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Poster not found"})
		return
	}

	c.Data(http.StatusOK, "image/jpeg", movie.Poster)
}

func (mh *MovieHandler) CreateMovie(c *gin.Context) {
	year, _ := strconv.Atoi(c.PostForm("year"))
	rating, _ := strconv.ParseFloat(c.PostForm("rating"), 64)

	movie := model.Movie{
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Year:        year,
		Rating:      rating,
	}

	file, err := c.FormFile("poster")
	if err == nil {
		openedFile, _ := file.Open()
		defer openedFile.Close()

		imageBytes, _ := io.ReadAll(openedFile)
		movie.Poster = imageBytes
	}

	createdMovie, err := mh.movieUsecase.CreateMovie(&movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
		return
	}

	c.JSON(http.StatusCreated, createdMovie)
}

// TODO: Implement UpdateMovie and DeleteMovie handlers
