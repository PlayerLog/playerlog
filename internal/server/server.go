package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/PlayerLog/playerlog/internal/auth"
	"github.com/PlayerLog/playerlog/internal/database"
	"github.com/PlayerLog/playerlog/internal/organization"
)

type Server struct {
	port int

	db                  database.Service
	authHandler         *auth.Handler
	organizationHandler *organization.Handler
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.New()
	NewServer := &Server{
		port:                port,
		db:                  db,
		authHandler:         auth.NewHandler(db, []byte("secret-key")),
		organizationHandler: organization.NewHandler(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
