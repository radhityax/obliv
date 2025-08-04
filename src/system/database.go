package system

import (
	"fmt"
	"database/sql"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

type Account struct {
	ID int
	Username string
	Password string
}

func CreateFile() {
	filename := "./data/obliv.db"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error while creating a file: %v\n", err)
		return
	}
	defer file.Close()
}

func ConnectDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./data/obliv.db")
	if err != nil {
		fmt.Printf("error while connecting database\n")
	}
	return db
}

func SetupDatabase(db *sql.DB) {
	stmt :=
	`
	CREATE TABLE IF NOT EXISTS account (
		USERNAME TEXT NOT NULL,
		PASSWORD TEXT NOT NULL
	)
	`

	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Printf("error while setting up database")
	}
}

func Register(db *sql.DB, username, password string) {
	stmt :=
	`
	INSERT INTO account (USERNAME, PASSWORD)
	VALUES (?,?)
	`
	_, err := db.Exec(stmt, username, password)
	if err != nil {
		fmt.Printf("error while inserting user")
	}
}
