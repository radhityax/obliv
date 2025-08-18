package system

import (
	"fmt"
_	"net/http"
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
	"time"
	"syscall"
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
	LoadAvg float64
	Idle uint64
	User uint64
	Nice uint64
	System uint64
	Iowait uint64
	Irq uint64
	Softirq uint64
	Steal uint64
	Guest uint64
	GuestNice uint64
	LoadOne float64
	LoadTwo float64
	LoadThree float64
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

func GetCpuStats() (CPU, error) {
	dat, err := ioutil.ReadFile("/proc/stat")
	check(err)

	lines := strings.Split(string(dat), "\n")
	cpuLine := lines[0]
	fields := strings.Fields(cpuLine)

	stats := CPU{}
	stats.User, _ = strconv.ParseUint(fields[1], 10, 64)
	stats.Nice, _ = strconv.ParseUint(fields[2], 10, 64)
	stats.System, _ = strconv.ParseUint(fields[3], 10, 64)
	stats.Idle,  _ = strconv.ParseUint(fields[4], 10, 64)
	stats.Iowait, _ = strconv.ParseUint(fields[5], 10, 64)
	stats.Irq, _ = strconv.ParseUint(fields[6], 10, 64)
	stats.Softirq, _ = strconv.ParseUint(fields[7], 10, 64)
	stats.Steal, _ = strconv.ParseUint(fields[8], 10, 64)
	stats.Guest, _ = strconv.ParseUint(fields[9], 10, 64)
	stats.GuestNice, _ = strconv.ParseUint(fields[10], 10, 64)
	return stats, err
}

func PrintCpu() (cpuUsage float64) {
	stats1, err := GetCpuStats()
	check(err)

	time.Sleep(1 * time.Second)

	stats2, err := GetCpuStats()
	check(err)

	totalTimeDiff := (
		stats2.User + stats2.Nice + stats2.System + stats2.Idle +
		stats2.Iowait + stats2.Irq + stats2.Softirq + stats2.Steal +
		stats2.Guest + stats2.GuestNice) - (stats1.User + stats1.Nice +
		stats1.System + stats1.Idle + stats1.Iowait +
		stats1.Irq + stats1.Softirq + stats1.Steal + stats1.Guest + stats1.GuestNice)

	idleTimeDiff := stats2.Idle - stats1.Idle

	if totalTimeDiff > 0 {
		cpuUsage = float64(totalTimeDiff-idleTimeDiff) / float64(totalTimeDiff) * 100
	}

	return cpuUsage
}

func GetLoadAverage() (float64, float64, float64, error) {
	file, err := os.Open("/proc/loadavg")
	if err != nil {
		return 0, 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 {
			return 0, 0, 0, fmt.Errorf("invalid format")
		}

		load1, err := parseLoad(fields[0])
		if err != nil {
			return 0, 0, 0, err
		}

		load5, err := parseLoad(fields[1])
		if err != nil {
			return 0, 0, 0, err
		}

		load15, err := parseLoad(fields[2])
		if err != nil {
			return 0, 0, 0, err
		}

		return load1, load5, load15, nil
	}

	return 0, 0, 0, fmt.Errorf("failed")
}

func parseLoad(loadStr string) (float64, error) {
	return strconv.ParseFloat(loadStr, 64)
}

func DiskUsage() (int, int, int, int) {
	var stat syscall.Statfs_t

	err := syscall.Statfs("/", &stat)
	if err != nil {
		fmt.Println("Error: ", err)
		return 0,0,0,0
	}

	blockSize := uint64(stat.Bsize)

	total := stat.Blocks * blockSize
	free := stat.Bfree * blockSize
	available := stat.Bavail * blockSize
	used := total - free

	const bytesInMB = 1024 * 1024
	totalMB := total / bytesInMB
	usedMB := used / bytesInMB
	freeMB := free / bytesInMB
	availableMB := available / bytesInMB

	return int(totalMB), int(usedMB), int(freeMB), int(availableMB)

}
