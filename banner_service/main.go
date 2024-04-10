package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var db *pgx.Conn

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	pgHost := os.Getenv("POSTGRES_CONTAINER_NAME")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDB := os.Getenv("POSTGRES_DB")

	pgConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pgUser, pgPassword, pgHost, pgPort, pgDB)
	conn, err := pgx.Connect(context.Background(), pgConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	db = conn

	http.HandleFunc("/", handleRequest)

	fmt.Println("Server is listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(context.Background(), `
		SELECT b.data_id, b.feature_id, bt.tag_id, b.is_active, b.created_at, b.updated_at
		FROM banners b
		INNER JOIN banner_tags bt ON b.id = bt.banner_id
	`)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying database: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	fmt.Fprintf(w, "| %10s | %10s | %10s | %10s | %30s | %30s |\n", "Data ID", "Feature ID", "Tag ID", "Is Active", "Created At", "Updated At")
	fmt.Fprintf(w, "|%s|\n", "------------+------------+------------+------------+------------------------------+------------------------------")
	for rows.Next() {
		var dataID string
		var featureID int64
		var tagID int64
		var isActive bool
		var createdAt time.Time
		var updatedAt time.Time

		err := rows.Scan(&dataID, &featureID, &tagID, &isActive, &createdAt, &updatedAt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "| %10s | %10d | %10d | %10t | %30s | %30s |\n", dataID, featureID, tagID, isActive, createdAt.Format(time.RFC3339), updatedAt.Format(time.RFC3339))
	}

	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error iterating over rows: %v", err), http.StatusInternalServerError)
		return
	}
}
