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
_	"obliv/src/system"
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

	http.ListenAndServe(Port, nil)
}
