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

func (r *AccountRepository) SaveCollection(user models.User, movieID int, collection string) (bool, error) {
	// Validate inputs
	if movieID <= 0 {
		r.logger.Error("SaveCollection failed: invalid movie ID", nil)
		return false, errors.New("invalid movie ID")
	}
	if collection != "favorite" && collection != "watchlist" {
		r.logger.Error("SaveCollection failed: invalid collection type", nil)
		return false, errors.New("collection must be 'favorite' or 'watchlist'")
	}

	// Get user ID from email
	var userID int
	err := r.db.QueryRow(`
		SELECT id 
		FROM users 
		WHERE email = $1 AND time_deleted IS NULL
	`, user.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found", nil)
		return false, ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("Failed to query user ID", err)
		return false, err
	}

	// Check if the relationship already exists
	var exists bool
	err = r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM user_movies 
			WHERE user_id = $1 
			AND movie_id = $2 
			AND relation_type = $3
		)
	`, userID, movieID, collection).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check existing collection entry", err)
		return false, err
	}
	if exists {
		r.logger.Info("Movie already in " + collection + " for user")
		return true, nil // Return true since the movie is already in the collection
	}

	// Insert the new relationship
	query := `
		INSERT INTO user_movies (user_id, movie_id, relation_type, time_added)
		VALUES ($1, $2, $3, $4)
	`
	_, err = r.db.Exec(query, userID, movieID, collection, time.Now())
	if err != nil {
		r.logger.Error("Failed to save movie to "+collection, err)
		return false, err
	}

	r.logger.Info("Successfully added movie " + string(movieID) + " to " + collection + " for user")
	return true, nil
}

var (
	ErrRegistrationValidation   = errors.New("registration failed")
	ErrAuthenticationValidation = errors.New("authentication failed")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrUserNotFound             = errors.New("user not found")
)
