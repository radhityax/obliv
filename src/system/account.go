package system

import (
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 
	bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ResetPassword(username, oldpass, newpass string) (string, error) {
	db := ConnectDatabase()
	defer db.Close()
	
	var existshash string
	
	err := db.QueryRow("SELECT password FROM account where username = ?",
	username).Scan(&existshash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("error no rows")
		}
		return "", fmt.Errorf("gagal")
		
	}

	if !CheckPassword(existshash, oldpass) {
		return "", fmt.Errorf("salah password lama ygy")
	}

	newpass, err = HashPassword(newpass)
	if err != nil {
		return "", fmt.Errorf("gagal hash password baru")
	}

	_, err = db.Exec("UPDATE account SET password = ? WHERE username = ?",
	newpass,username)
	if err != nil {
		return "", fmt.Errorf("gagal update password baru")
	}

	return newpass, nil
}
func Register(db *sql.DB, username, password string) error {
        if len(username) < 3 || len(password) < 8 {
                return fmt.Errorf("username < 3 || password < 8")
        }

        hashpass, err := HashPassword(password)

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
