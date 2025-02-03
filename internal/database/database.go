package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	//_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/stephenafamo/bob"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string
	GetDB() bob.DB
	// Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	db bob.DB
}

var (
	database   = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	schema     = os.Getenv("BLUEPRINT_DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	// db, err := sql.Open("pgx", connStr)
	db, err := bob.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// 	_, err = db.Exec(`

	// CREATE TABLE teams (
	//     id TEXT PRIMARY KEY,  -- CUID
	//     name VARCHAR(100) NOT NULL,
	//     short_name VARCHAR(20),
	//     logo_url TEXT,
	//     primary_color VARCHAR(7),  -- Hex color code
	//     secondary_color VARCHAR(7),  -- Hex color code
	//     sport_type VARCHAR(50) NOT NULL,  -- e.g., 'football', 'hockey', etc.
	//     --club_id TEXT REFERENCES clubs(id), -- if team belongs to a larger club
	//     division VARCHAR(50),
	//     season VARCHAR(20),
	//     is_active BOOLEAN DEFAULT true,
	//     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	//     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	// );

	// -- For team settings/preferences
	// CREATE TABLE team_settings (
	//     team_id TEXT PRIMARY KEY REFERENCES teams(id) ON DELETE CASCADE,
	//     privacy_level VARCHAR(20) DEFAULT 'private',  -- 'private', 'public', 'unlisted'
	//     stats_visibility JSONB DEFAULT '{"public": false, "members_only": true}',
	//     match_reminder_hours INTEGER DEFAULT 24,
	//     default_match_duration INTEGER DEFAULT 90  -- in minutes
	// );

	// -- Indexes
	// --CREATE INDEX idx_teams_club_id ON teams(club_id);
	// CREATE INDEX idx_teams_sport_type ON teams(sport_type);

	// 	CREATE TABLE users (
	//     id TEXT PRIMARY KEY,  -- CUID
	//     email VARCHAR(255) UNIQUE NOT NULL,
	//     password_hash VARCHAR(255) NOT NULL,
	//     first_name VARCHAR(100),
	//     last_name VARCHAR(100),
	//     role VARCHAR(50) NOT NULL, -- 'admin', 'team_admin', 'coach', 'player', 'parent'
	//     phone VARCHAR(20),
	//     is_active BOOLEAN DEFAULT true,
	//     last_login TIMESTAMP,
	//     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	//     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	// );

	// CREATE TABLE user_settings (
	//     user_id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
	//     notification_preferences JSONB DEFAULT '{"email": true, "push": false}',
	//     theme VARCHAR(20) DEFAULT 'light',
	//     language VARCHAR(10) DEFAULT 'en',
	//     timezone VARCHAR(50) DEFAULT 'UTC'
	// );

	// CREATE TABLE team_members (
	//     user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
	//     team_id TEXT REFERENCES teams(id) ON DELETE CASCADE,
	//     role VARCHAR(50) NOT NULL,
	//     jersey_number INTEGER,
	//     position VARCHAR(50),
	//     joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	//     PRIMARY KEY (user_id, team_id)
	// );

	// --CREATE TABLE sessions (
	// --    id TEXT PRIMARY KEY,
	// --    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
	// --    token TEXT UNIQUE NOT NULL,
	// --    expires_at TIMESTAMP NOT NULL,
	// --    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	// --    ip_address VARCHAR(45),
	// --    user_agent TEXT
	// --);

	// CREATE INDEX idx_users_email ON users(email);
	// CREATE INDEX idx_team_members_team_id ON team_members(team_id);
	// --CREATE INDEX idx_sessions_token ON sessions(token);
	// --CREATE INDEX idx_sessions_user_id ON sessions(user_id);
	// 	`)

	fmt.Println(err)

	dbInstance = &service{
		db: db,
	}

	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	_, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	// err := s.db.PingContext(ctx)
	// if err != nil {
	// 	stats["status"] = "down"
	// 	stats["error"] = fmt.Sprintf("db down: %v", err)
	// 	log.Fatalf("db down: %v", err) // Log the error and terminate the program
	// 	return stats
	// }

	// // Database is up, add more statistics
	// stats["status"] = "up"
	// stats["message"] = "It's healthy"

	// // Get database stats (like open connections, in use, idle, etc.)
	// dbStats := s.db.Stats()
	// stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	// stats["in_use"] = strconv.Itoa(dbStats.InUse)
	// stats["idle"] = strconv.Itoa(dbStats.Idle)
	// stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	// stats["wait_duration"] = dbStats.WaitDuration.String()
	// stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	// stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// // Evaluate stats to provide a health message
	// if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
	// 	stats["message"] = "The database is experiencing heavy load."
	// }

	// if dbStats.WaitCount > 1000 {
	// 	stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	// }

	// if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
	// 	stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	// }

	// if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
	// 	stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	// }

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

// func (s *service) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
// 	return s.db.QueryContext(ctx, query, args...)
// }

// func (s *service) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
// 	return s.db.ExecContext(ctx, query, args...)
// }

func (s *service) GetDB() bob.DB {
	return s.db
}
