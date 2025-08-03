package system

import (
	"fmt"
	"database/sql"
	"os"
)

type Account struct {
	ID int
	Username string
	Password string
}

func createfile() {
	filename := "./data/obliv.db"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("error while creating a file: %v\n", err)
		return
	}
	defer file.Close()
}

func connectdatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./data/obliv.db")
	if err != nil {
		fmt.Printf("error while connecting database\n")
	}
	return db
}
