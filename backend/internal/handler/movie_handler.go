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
	year, err := strconv.Atoi(c.PostForm("year"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format"})
		return
	}

	rating, err := strconv.ParseFloat(c.PostForm("rating"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating format"})
		return
	}

	movie := model.Movie{
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Year:        year,
		Rating:      rating,
	}

	file, err := c.FormFile("poster")
	if err != nil {
		fmt.Printf("Error getting poster file: %v\n", err)
	} else {
		fmt.Printf("File received: %s, Size: %d bytes\n", file.Filename, file.Size)

		openedFile, err := file.Open()
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open poster file"})
			return
		}
		defer openedFile.Close()

		imageBytes, err := io.ReadAll(openedFile)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read poster file"})
			return
		}

		fmt.Printf("Image bytes read: %d bytes\n", len(imageBytes))
		movie.Poster = imageBytes
	}

	createdMovie, err := mh.movieUsecase.CreateMovie(&movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
		return
	}

	createdMovie.PosterURL = fmt.Sprintf("/api/v1/movies/%d/poster", createdMovie.ID)
	c.JSON(http.StatusCreated, createdMovie)
}

// TODO: Implement UpdateMovie and DeleteMovie handlers
