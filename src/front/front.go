package front

import (
	"net/http"
	"fmt"
	"encoding/json"
	"obliv/src/system"
)

var head = `
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>obliv</title>
</head>
`

var footer = `
</html>
`

func Homepage(w http.ResponseWriter, r *http.Request) {
	html := `
	<body>
	<p>this is obliv homepage</p>
	</body>
	`
	fmt.Fprintf(w, head)
	fmt.Fprintf(w, html)
}
func MemoryPage(w http.ResponseWriter, r *http.Request) {
	dat := system.GetMemory()
	html := `
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
	fmt.Fprintf(w, head)
	fmt.Fprintf(w, html, dat.MemTotal,
	dat.MemTotal-dat.MemFree-dat.Buffers-dat.Cached, dat.MemAvailable, 
	dat.Buffers, dat.Cached)
	fmt.Fprintf(w, footer)
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

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		dtb := system.ConnectDatabase()
		defer dtb.Close()

		system.Register(dtb, username, password)

		// debug
		fmt.Printf("User submitted : %s\n", username)
		fmt.Printf("Password submitted : %s\n", password)
		return
	}

	body :=
	`
	<p>register</p>
	<form method="POST">
	<input type="text" name="username" placeholder="username">
	<input type="password" name="password" placeholder="password">
	<button type="submit">register</button>
	</form>
	`
	fmt.Fprintf(w, head)
	fmt.Fprintf(w, body)
	fmt.Fprintf(w, footer)
}

func Loginpage(w http.ResponseWriter, r *http.Request) {
}
