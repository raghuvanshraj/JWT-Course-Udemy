package driver

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
	"os"
)

func ConnectDB() *sql.DB {
	logger := log.New(os.Stdout, "db-connection: ", log.LstdFlags)

	pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	if err != nil {
		logger.Fatalln(err)
	}

	db, err := sql.Open("postgres", pgUrl)
	if err != nil {
		logger.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Println("connection to db successful")

	return db
}
