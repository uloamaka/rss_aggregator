package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"context"
	"os/signal"
    "syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	database "github.com/uloamaka/rss_aggregator/internal/database"
	"github.com/jackc/pgx/v5"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	ctx := context.Background()

	godotenv.Load(".env")

	
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}
	
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	// conn, err := sql.Open("postgres", dbURL)"user=pqgotest dbname=pqgotest sslmode=verify-full"
	// if err != nil {
	// 	log.Fatal("Can't connect to database:", err)
	// }

	conn, err := pgx.Connect(ctx, dbURL )
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	defer conn.Close(ctx)

	// queries := tutorial.New(conn)

	apiCfg := apiConfig {
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerError)
	v1Router.Post("/user", apiCfg.handlerCreateUser)
	v1Router.Get("/user", apiCfg.handlerGetUser)


	router.Mount("/api/v1", v1Router)

	srv := &http.Server {
		Handler: router,
		Addr: ":" + portString,
	}

	go func() {
        fmt.Println("Server running on port", portString)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("ListenAndServe(): %s", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit

    fmt.Println("Shutting down server...")

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server Shutdown Failed:%+v", err)
    }
    fmt.Println("Server exited properly")

}