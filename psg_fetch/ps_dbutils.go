package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

func main() {
	db, err := SetupDB()
	if err != nil {
		log.Fatal(err)
	}

	affect, err := remove_Dups(db, "artist", "artistid")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Done updating")
}

func remove_Dups(db *sql.DB, string tablename, string column) (int, error) {
	queryText :=
		`DELETE FROM $1 
		WHERE id IN (SELECT id
				FROM (SELECT id,
									ROW_NUMBER() OVER 
										(partition BY $2 ORDER BY id) AS rnum
							FROM $1) t
				WHERE t.rnum > 1);`

	log.Println("Removing dups from :", tablename)
	queryUpdate, err := db.Prepare(queryText)
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.Exec(tablename, column)
	if err != nil {
		log.Fatal(err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Duped rows affected :", affect)
	return affect, err
}

func SetupDB() (*sql.DB, error) {
	username := os.Getenv("DBUSER_WRITE") // "nick" // for dev
	log.Println("Username: ", username)
	pass := os.Getenv("DBUSER_WRITE_PW") // "nick" // for dev
	database := os.Getenv("PLAYLISTDB")
	app_status := os.Getenv("APP_STATUS")
	var connection string
	if app_status == "DEV" {
		connection = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", username, pass, database)
	} else {
		connection = fmt.Sprintf("user=%s password=%s dbname=%s", username, pass, database)
	}

	log.Println("Opening database:", database)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal(err)
	}

	return db, err

}
