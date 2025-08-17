package system

import (
	"fmt"
	"os"
	"bufio"
	"strings"
)

func Shell() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("obliv> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		input = strings.TrimSpace(input)
		ShellMenu(input)
	}
}

func ShellMenu(input string) {
	switch(input) {
	case "H", "h":
		fmt.Println("n - new user")
		fmt.Println("h - help")
		fmt.Println("r - reset password")
		fmt.Println("v - version")
	case "R", "r":
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Input username: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("failed to read username: ", err)
			return
		}
		username = strings.TrimSpace(username)

		fmt.Printf("Input oldpass: ")
	oldpass, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("failed to read oldpass: ", err)
			return
		}
		oldpass = strings.TrimSpace(oldpass)

		fmt.Printf("Input newpass: ")
		newpass, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("failed to read newpass: ", err)
			return
		}
		newpass = strings.TrimSpace(newpass)
		
		_, err = ResetPassword(username, oldpass, newpass) 
	case "N", "n":
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Input username: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("failed to read username:", err)
			return
		}
		username = strings.TrimSpace(username)
		
		fmt.Printf("Input password: ")
		password, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("failed to read password: ", err)
		}
		password = strings.TrimSpace(password)
			
		db := ConnectDatabase()
		defer db.Close()
		err = Register(db, username, password)
		if err != nil {
			fmt.Println("Error while creating an account: ", err)
			return
		}
		fmt.Println("akun berhasil dibuat")
	default:
		fmt.Println("perintah apa tuh? coba ketik 'h' aja")
	}
}
