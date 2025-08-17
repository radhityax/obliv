package system

import (
	"fmt"
	"database/sql"
	"os"
	_ "github.com/mattn/go-sqlite3"
	_ "golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID int
	Username string
	Password string
}

func CreateFile() error {
	os.Mkdir("./data", os.ModeDir)
	filename := "./data/obliv.db"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		file.Close()
	}
	return nil
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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)
	`

	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Printf("error while setting up database: %v\n", err)
	}
}
