package front

import (
	"net/http"
	"fmt"
	"encoding/json"
	"obliv/src/system"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	html := `
	<html>
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>obliv</title>
	</head>
	<body>
	<p>this is obliv</p>
	</body>
	</html>
	`
	fmt.Fprintf(w, html)
}
func MemoryPage(w http.ResponseWriter, r *http.Request) {
	dat := system.GetMemory()
	html := `
	<html>
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>obliv</title>
	<script>
	function updateMemory() {
		fetch('/memory-data')
		.then(response => response.json())
		.then(data => {
			document.getElementById("memTotal").innerText = 'Total: ' + data.memTotal;
			document.getElementById("memUsed").innerText = 'Used: ' + data.memUsed;
			document.getElementById("memAvailable").innerText = 'Available: ' + data.memAvailable;
			document.getElementById("buffers").innerText = 'Buffers: ' + data.buffers;
			document.getElementById("cached").innerText = 'Cached: ' + data.cached;
		});
	}
	setInterval(updateMemory, 3000);
	</script>
	</head>
	<body>
	<p>this is memory</p>
	<p id="memTotal">Total: %d</p>
	<p id="memUsed">Used: %d</p>
	<p id="memAvailable">Available: %d</p>
	<p id="buffers">Buffers: %d</p>
	<p id="cached">Cached: %d</p>
	</body>
	</html>
	`
	fmt.Fprintf(w, html, dat.MemTotal,
	dat.MemTotal-dat.MemFree-dat.Buffers-dat.Cached, dat.MemAvailable, 
	dat.Buffers, dat.Cached)
}

func MemoryData(w http.ResponseWriter, r *http.Request) {
	dat := system.GetMemory()
	response := map[string]interface{}{
		"memTotal":    dat.MemTotal,
		"memUsed":     dat.MemTotal - dat.MemFree - dat.Buffers - dat.Cached,
		"memAvailable": dat.MemAvailable,
		"buffers":     dat.Buffers,
		"cached":      dat.Cached,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func CpuPage(w http.ResponseWriter, r *http.Request) {
	html := `
	<html>
	<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="refresh" content="1" />
	<title>obliv</title>
	</head>
	<body>
	<p>this is cpu</p>
	<p>average: %f</p>
	</body>
	</html>
	`

	fmt.Fprintf(w, html, system.PrintCPU())
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
}

func Loginpage(w http.ResponseWriter, r *http.Request) {
}
