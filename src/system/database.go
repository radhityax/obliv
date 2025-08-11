package system

import (
	"fmt"
	"database/sql"
	"os"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID int
	Username string
	Password string
}

func CreateFile() error {
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

func Register(db *sql.DB, username, password string) error {
	if len(username) < 3 || len(password) < 8 {
		return fmt.Errorf("username < 3 || password < 8")
	}

	hashpass, err := bcrypt.GenerateFromPassword([]byte(password),
	bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate a password: %v", err)
	}

	stmt := `
	INSERT INTO account (USERNAME, PASSWORD)
	VALUES (?,?)
	`

	result, err := db.Exec(stmt, username, hashpass)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: account.username" {
			return fmt.Errorf("username %s sudah digunakan", err)
		}
		return fmt.Errorf("gagal menyimpan pengguna: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("tidak ada data yang disimpan")
	}

	return nil
}
