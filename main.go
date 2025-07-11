package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mohndakbar/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int64
	dbQueires      *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dnQueries := database.New(db)

	const filePathRoot = "."
	const port = "8080"
	apiCfg := &apiConfig{
		dbQueires: dnQueries,
		platform:  platform,
		jwtSecret: jwtSecret,
		polkaKey:  polkaKey,
	}

	// ServeMux is an HTTP request multiplexer.
	// It matches the URL of each incoming request against a list of registered patterns
	// and calls the handler for the pattern that most closely matches the URL.
	mux := http.NewServeMux()

	// A Handler responds to an HTTP request.
	// http.FileServer is a Handler that serves HTTP requests with the contents of the file system root directory.
	mux.Handle("/app/", apiCfg.middlewareMetricIncrement(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot))))) // serve static files(http.FileServer(http.Dir(".")))
	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerDeleteAllUsers)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirp_id}", apiCfg.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirp_id}", apiCfg.handlerDeleteChirp)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefreshToekn)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerCreateWebhook)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

// func main() {
// 	// main impelementation with gorilla/mux library
// 	const filePathRoot = "."
// 	const port = "8080"
// 	apiCfg := &apiConfig{}
// 	// ServeMux is an HTTP request multiplexer.
// 	// It matches the URL of each incoming request against a list of registered patterns
// 	// and calls the handler for the pattern that most closely matches the URL.
// 	router := mux.NewRouter()
// 	// A Handler responds to an HTTP request.
// 	// http.FileServer is a Handler that serves HTTP requests with the contents of the file system root directory.
// 	router.Handle("/app/", apiCfg.middlewareMetricIncrement(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot))))) // serve static files(http.FileServer(http.Dir(".")))
// 	router.HandleFunc("/healthz", handlerHealth).Methods("GET")
// 	router.HandleFunc("/metrics", apiCfg.handlerMetrics).Methods("GET")
// 	router.HandleFunc("/reset", apiCfg.handlerResetMetrics).Methods("POST")
// 	srv := &http.Server{
// 		Addr:    ":" + port,
// 		Handler: router,
// 	}
// 	log.Printf("Serving files from %s on port %s", filePathRoot, port)
// 	log.Fatal(srv.ListenAndServe())
// }
