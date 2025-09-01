package data

import (
	"SagiProjects.com/moviesite/models"
)

type MovieStorage interface {
	GetTopMovies() ([]models.Movie, error)
	GetRandomMovies() ([]models.Movie, error)
	GetMovieByID(id int) (models.Movie, error)
	SearchMoviesByName(name string, order string, genre *int) ([]models.Movie, error)
	GetAllGenres() ([]models.Genre, error)
}

type AccountStorage interface {
	Authenticate(string, string) (*models.User, error)
	Register(string, string, string) (bool, error)
	GetAccountDetails(string) (models.User, error)
	SaveCollection(models.User, int, string) (models.User, error)
}
