package main

import (
	"fmt"
	"net/http"
	"bufio"
_	"io"
	"os"
	"strings"
	"strconv"

	"obliv/src/front"
)

var Port string = ":2305"

type Memory struct {
	MemTotal int
	MemFree int
	MemAvailable int
	Buffers int
	Cached int
	BuffCache int
}
type CPU struct {
	Usage float32
}
func GetMemory() Memory {
	dat, err := os.Open("/proc/meminfo")
	check(err)
	defer dat.Close()
	
	scan := bufio.NewScanner(dat)
	res := Memory {}
	for scan.Scan() {
		key, value := parseLine(scan.Text())
		switch key {
			case "MemTotal":
				res.MemTotal = value / 1024
			case "MemFree":
				res.MemFree = value / 1024
			case "MemAvailable":
				res.MemAvailable = value / 1024
			case "Buffers":
				res.Buffers = value / 1204
			case "Cached":
				res.Cached = value / 1204
			}
		}
		return res
}
func parseLine(raw string) (key string, value int) {
    text := strings.ReplaceAll(raw[:len(raw)-2], " ", "")
    keyValue := strings.Split(text, ":")
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
	fmt.Fprintf(w, "Total: %d\n", dat.MemTotal)
	fmt.Fprintf(w, "Used: %d\n", 
		dat.MemTotal - dat.MemFree - dat.Buffers - dat.Cached)
	fmt.Fprintf(w, "Available: %d\n", dat.MemAvailable)
	fmt.Fprintf(w, "Buff/Cache: %d/%d\n", dat.Buffers,dat.Cached)
}

// func GetCPU() CPU {
//}

func PrintCPU(w http.ResponseWriter, r *http.Request) {
}

func main() {
	intro()

	http.HandleFunc("/", front.Homepage)
	http.HandleFunc("/ping", pong)
	http.HandleFunc("/memory", PrintMemory)
	http.ListenAndServe(Port, nil)
}
