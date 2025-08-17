package main

import (
	"fmt"
	"net/http"

_	"bufio"
_	"io/ioutil"
_	"os"
_	"strings"
_	"strconv"
_	"time"

	"obliv/src/front"
	"obliv/src/system"
)

var Port string = ":2305"

func intro() {
	fmt.Println("obliv - control panel [BETA]")
	fmt.Println("https://github.com/radh1tya/obliv");
}

func pong(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong!")
}

func main() {
	intro()
	
	if err := system.CreateFile(); err != nil {
		fmt.Printf("Failed while creating a database file: %v\n", err)
		return
	}

	dtb := system.ConnectDatabase()

	defer dtb.Close()

	system.SetupDatabase(dtb)

	http.HandleFunc("/", front.Homepage)
	http.HandleFunc("/ping", pong)
	http.HandleFunc("/memory", front.MemoryPage)
	http.HandleFunc("/memory-data", front.MemoryData)
	
	http.HandleFunc("/register", front.RegisterPage)
	http.HandleFunc("/login", front.Loginpage)
	
	/*
	i dont know the best way to know about it
	http.HandleFunc("/cpu", front.CpuPage)
	*/

	fmt.Printf("currently running on %s\n", Port)
	if err := http.ListenAndServe(Port, nil); err != nil {
		fmt.Printf("failed run server: %v\n", err)
	}
}
