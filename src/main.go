package main

import (
	"fmt"
	"net/http"
	"bufio"
	_ "io"
	"os"
	"strings"
	"strconv"
)

type Memory struct {
	MemTotal int
	MemFree int
	MemAvailable int
}

func GetMemory() Memory {
	dat, err := os.Open("/proc/meminfo")
	if err != nil {
		panic(err)
	}
	defer dat.Close()
	
	scan := bufio.NewScanner(dat)
	res := Memory {}
	for scan.Scan() {
		key, value := parseLine(scan.Text())
		switch key {
			case "MemTotal":
				res.MemTotal = value / 1024
			case "MemFree":
				res.MemFree = value
			case "MemAvailable":
				res.MemAvailable = value
			}
		}
		return res
}
func parseLine(raw string) (key string, value int) {
    text := strings.ReplaceAll(raw[:len(raw)-2], " ", "")
    keyValue := strings.Split(text, ":")
    fmt.Println(keyValue[1])
	return keyValue[0], toInt(keyValue[1])
}

func toInt(raw string) int {
    if raw == "" {
        return 0
    }
    res, err := strconv.Atoi(raw)
    
	check(err)
	return res
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func intro() {
	fmt.Println("obliv - control panel [BETA]")
	fmt.Println("https://github.com/radh1tya/obliv");
}

func pong(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Pong!")
}

func PrintMemory(w http.ResponseWriter, r *http.Request) {
	dat := GetMemory()
	fmt.Fprintf(w, "Currently using: %d", dat.MemTotal)
}

func main() {
	intro()

	http.HandleFunc("/ping", pong)
	http.HandleFunc("/memory", PrintMemory)
	http.ListenAndServe(":8080", nil)
}
