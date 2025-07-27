package system

import (
_	"fmt"
_	"net/http"
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
	"time"
)

type Memory struct {
	MemTotal int
	MemFree int
	MemAvailable int
	Buffers int
	Cached int
	BuffCache int
}
type CPU struct {
	loadAvg float64
	idle uint64
	user uint64
	nice uint64
	system uint64
	iowait uint64
	irq uint64
	softirq uint64
	steal uint64
	guest uint64
	guestNice uint64
}
func check(e error) {
	if e != nil {
		panic(e)
	}
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
/*
func PrintMemory(w http.ResponseWriter, r *http.Request) {
	dat := GetMemory()
	fmt.Fprintf(w, "Total: %d\n", dat.MemTotal)
	fmt.Fprintf(w, "Used: %d\n",
		dat.MemTotal - dat.MemFree - dat.Buffers - dat.Cached)
	fmt.Fprintf(w, "Available: %d\n", dat.MemAvailable)
	fmt.Fprintf(w, "Buff/Cache: %d/%d\n", dat.Buffers,dat.Cached)
}
*/


func getCpuStats() (CPU, error) {
	dat, err := ioutil.ReadFile("/proc/stat")
	check(err)

	lines := strings.Split(string(dat), "\n")
	cpuLine := lines[0]
	fields := strings.Fields(cpuLine)

	stats := CPU{}
	stats.user, _ = strconv.ParseUint(fields[1], 10, 64)
	stats.nice, _ = strconv.ParseUint(fields[2], 10, 64)
	stats.system, _ = strconv.ParseUint(fields[3], 10, 64)
	stats.idle,  _ = strconv.ParseUint(fields[4], 10, 64)
	stats.iowait, _ = strconv.ParseUint(fields[5], 10, 64)
	stats.irq, _ = strconv.ParseUint(fields[6], 10, 64)
	stats.softirq, _ = strconv.ParseUint(fields[7], 10, 64)
	stats.steal, _ = strconv.ParseUint(fields[8], 10, 64)
	stats.guest, _ = strconv.ParseUint(fields[9], 10, 64)
	stats.guestNice, _ = strconv.ParseUint(fields[10], 10, 64)
	return stats, err
}

func PrintCPU() (cpuUsage float64) {
	stats1, err := getCpuStats()
	check(err)

	time.Sleep(1 * time.Second)

	stats2, err := getCpuStats()
	check(err)

	totalTimeDiff := (
		stats2.user + stats2.nice + stats2.system + stats2.idle +
		stats2.iowait + stats2.irq + stats2.softirq + stats2.steal +
		stats2.guest + stats2.guestNice) - (stats1.user + stats1.nice +
		stats1.system + stats1.idle + stats1.iowait +
		stats1.irq + stats1.softirq + stats1.steal + stats1.guest + stats1.guestNice)

	idleTimeDiff := stats2.idle - stats1.idle

	if totalTimeDiff > 0 {
		cpuUsage = float64(totalTimeDiff-idleTimeDiff) / float64(totalTimeDiff) * 100
	}

	return cpuUsage
}
