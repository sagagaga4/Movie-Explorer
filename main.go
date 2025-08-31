package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	logger "SagiProjects.com/moviesite/Logger"
	"SagiProjects.com/moviesite/data"
	"SagiProjects.com/moviesite/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func initializeLogger() *logger.Logger {
	logInstance, err := logger.NewLogger("movie.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: $v", err)
	}
	defer logInstance.Close()
	return logInstance
}

func main() {

	//Log Initializer
	logInstance := initializeLogger()

	//Enviromental Variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file")
	}

	//Connect to the DB
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		log.Fatalf("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: $v", err)
	}
	defer db.Close()

	//Initialize Data Repository for Movies
	movieRepo, err := data.NewMovieRepository(db, logInstance)
	if err != nil {
		log.Fatalf("Failed to initialize movie repository: $v", err)
	}

	//Movie Handler Initializer using New to get the "private" varibles
	//Accessing the storage for passing the db
	movieHandler := handlers.NewMovieHandler(movieRepo, logInstance)

	//Initialize Account Repository for Users
	accountRepo, err := data.NewAccountRepository(db, logInstance)
	if err != nil {
		log.Fatalf("Failed to initialize account repository")
	}
	//Account Handler Initializer
	accountHandler := handlers.NewAccountHandler(accountRepo, logInstance)

	//Handlers Calls
	http.HandleFunc("/api/movies/top/", movieHandler.GetTopMovies)
	http.HandleFunc("/api/movies/random/", movieHandler.GetRandomMovies)
	http.HandleFunc("/api/movies/search/", movieHandler.SearchMovies)
	http.HandleFunc("/api/movies/", movieHandler.GetMovie)
	http.HandleFunc("/api/genres/", movieHandler.GetGenres)
	http.HandleFunc("/api/account/register/", accountHandler.Register)
	http.HandleFunc("/api/account/authenticate/", accountHandler.Authenticate)

	catchAllClientRoutesHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	}
	http.HandleFunc("/movies", catchAllClientRoutesHandler)
	http.HandleFunc("/movies/", catchAllClientRoutesHandler)
	http.HandleFunc("/account/", catchAllClientRoutesHandler)

	//Handler for static files(Front-End)
	http.Handle("/", http.FileServer(http.Dir("public")))
	fmt.Println("Serving the files ")

	const addr = "localhost:8080"
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
		logInstance.Error("Server failed", err)
	}
}
