package connection

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Create() *sql.DB {
	// Connection string
	connStr := "user=myuser password=helloworld dbname=analytics sslmode=disable"

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		defer db.Close()
		log.Fatal("Failed to connect to the database:", err)
	}
	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("Connected to the database successfully")
	return db
}
