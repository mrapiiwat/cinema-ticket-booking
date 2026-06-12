package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
)

type MovieHandler struct {
	movieService ports.MovieServicePort
}

func NewMovieHandler(movieService ports.MovieServicePort) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}

func (h *MovieHandler) GetMovies(c echo.Context) error {
	movies, err := h.movieService.GetMovies(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, movies)
}

func (h *MovieHandler) GetMovieByID(c echo.Context) error {
	id := c.Param("id")
	movie, err := h.movieService.GetMovieByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if movie == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Movie not found"})
	}
	return c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) CreateMovie(c echo.Context) error {
	var movie model.Movie
	if err := c.Bind(&movie); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.movieService.CreateMovie(c.Request().Context(), &movie); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, movie)
}
