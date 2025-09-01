package data

import (
	"database/sql"
	"errors"
	"time"

	logger "SagiProjects.com/moviesite/Logger"
	"SagiProjects.com/moviesite/models"
	"golang.org/x/crypto/bcrypt"
)

type AccountRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewAccountRepository(db *sql.DB, log *logger.Logger) (*AccountRepository, error) {
	return &AccountRepository{
		db:     db,
		logger: log,
	}, nil
}

func (r *AccountRepository) Register(name, email, password string) (bool, error) {
	// Validate basic requirements
	if name == "" || email == "" || password == "" {
		r.logger.Error("Registration validation failed: missing required fields", nil)
		return false, ErrRegistrationValidation
	}

	// Check if user already exists
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`, email).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check existing user", err)
		return false, err
	}
	if exists {
		r.logger.Error("User already exists with email: "+email, ErrUserAlreadyExists)
		return false, ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error("Failed to hash password", err)
		return false, err
	}

	// Insert new user
	query := `
		INSERT INTO users (name, email, password_hashed, time_created)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var userID int
	err = r.db.QueryRow(
		query,
		name,
		email,
		string(hashedPassword),
		time.Now(),
	).Scan(&userID)
	if err != nil {
		r.logger.Error("Failed to register user", err)
		return false, err
	}

	return true, nil
}

func (r *AccountRepository) Authenticate(email string, password string) (*models.User, error) {
	if email == "" || password == "" {
		r.logger.Error("Authentication validation failed: missing credentials", nil)
		return nil, ErrAuthenticationValidation
	}

	// Fetch user by email
	var user models.User
	query := `
		SELECT id, name, email, password_hashed
		FROM users 
		WHERE email = $1 AND time_deleted IS NULL
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found for email: "+email, nil)
		return nil, ErrAuthenticationValidation
	}
	if err != nil {
		r.logger.Error("Failed to query user for authentication", err)
		return nil, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		r.logger.Error("Password mismatch for email: "+email, nil)
		return nil, ErrAuthenticationValidation
	}

	// Update last login time
	updateQuery := `
		UPDATE users 
		SET last_login = $1
		WHERE id = $2
	`
	_, err = r.db.Exec(updateQuery, time.Now(), user.ID)
	if err != nil {
		r.logger.Error("Failed to update last login", err)
		// Don't fail authentication just because last login update failed
	}

	return &user, nil
}

func (r *AccountRepository) GetAccountDetails(email string) (models.User, error) {
	var user models.User
	query := `
		SELECT id, name, email
		FROM users 
		WHERE email = $1 AND time_deleted IS NULL
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found for email: "+email, nil)
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("Failed to query user by email", err)
		return models.User{}, err
	}

	// Fetch favorites
	favoritesQuery := `
		SELECT m.id, m.tmdb_id, m.title, m.tagline, m.release_year, 
		       m.overview, m.score, m.popularity, m.language, 
		       m.poster_url, m.trailer_url
		FROM movies m
		JOIN user_movies um ON m.id = um.movie_id
		WHERE um.user_id = $1 AND um.relation_type = 'favorite'
	`
	favoriteRows, err := r.db.Query(favoritesQuery, user.ID)
	if err != nil {
		r.logger.Error("Failed to query user favorites", err)
		return user, err
	}
	defer favoriteRows.Close()

	for favoriteRows.Next() {
		var m models.Movie
		if err := favoriteRows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan favorite movie row", err)
			return user, err
		}
		user.Favorites = append(user.Favorites, m)
	}

	// Fetch watchlist
	watchlistQuery := `
		SELECT m.id, m.tmdb_id, m.title, m.tagline, m.release_year, 
		       m.overview, m.score, m.popularity, m.language, 
		       m.poster_url, m.trailer_url
		FROM movies m
		JOIN user_movies um ON m.id = um.movie_id
		WHERE um.user_id = $1 AND um.relation_type = 'watchlist'
	`
	watchlistRows, err := r.db.Query(watchlistQuery, user.ID)
	if err != nil {
		r.logger.Error("Failed to query user watchlist", err)
		return user, err
	}
	defer watchlistRows.Close()

	for watchlistRows.Next() {
		var m models.Movie
		if err := watchlistRows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan watchlist movie row", err)
			return user, err
		}
		user.Watchlist = append(user.Watchlist, m)
	}

	return user, nil
}

// Change the return signature from (bool, error) to (models.User, error)
func (r *AccountRepository) SaveCollection(user models.User, movieID int, collectionType string) (models.User, error) {
	var query string
	if collectionType == "favorites" {
		query = `INSERT INTO user_favorites (user_id, movie_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	} else if collectionType == "watchlist" {
		query = `INSERT INTO user_watchlist (user_id, movie_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	} else {
		return user, errors.New("invalid collection type")
	}

	_, err := r.db.Exec(query, user.ID, movieID)
	if err != nil {
		r.logger.Error("Failed to save to collection", err)
		return user, err
	}

	// After a successful save, return the user object.
	// For a more robust implementation, you might refetch the user's collections here.
	// For now, we return the user object as is.
	return user, nil
}

var (
	ErrRegistrationValidation   = errors.New("registration failed")
	ErrAuthenticationValidation = errors.New("authentication failed")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrUserNotFound             = errors.New("user not found")
)
