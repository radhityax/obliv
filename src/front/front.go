package front

import (
	"net/http"
	"fmt"
	"encoding/json"
	"obliv/src/system"
	"golang.org/x/crypto/bcrypt"
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
		err := r.ParseForm()
		
		if err != nil {
			http.Error(w, "failed while processing a form",
			http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		db := system.ConnectDatabase()
		defer db.Close()

		err = system.Register(db, username, password)
		if err != nil {
			fmt.Fprintf(w, head)
			fmt.Fprintf(w, `
			<body>
			<p style="color: red;">Failed: %s</p>
			<p><a href="/register">Try again</a></p>
			</body>
			`, err.Error())
			fmt.Fprintf(w, footer)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed while processing form",
			http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		db := system.ConnectDatabase()
		defer db.Close()

		var hashedPassword string
		query := `SELECT password FROM account WHERE username = ?`
		err = db.QueryRow(query, username).Scan(&hashedPassword)
		if err != nil {
			fmt.Fprintf(w, head)
			fmt.Fprintf(w, `
			<body>
			<p style="color: red;">cant find the username</p>
			<p><a href="/login">try again</a></p>
			</body>
			`)
			fmt.Fprintf(w, footer)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
        if err != nil {
            fmt.Fprintf(w, head)
            fmt.Fprintf(w, `
            <body>
            <p style="color: red;">Password salah</p>
            <p><a href="/login">Coba lagi</a></p>
            </body>
            `)
            fmt.Fprintf(w, footer)
            return
        }

        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    body := `
    <body>
    <p>login</p>
    <form method="POST">
    <input type="text" name="username" placeholder="username" required>
    <input type="password" name="password" placeholder="password" required>
    <button type="submit">login</button>
    </form>
    </body>
    `
    fmt.Fprintf(w, head)
    fmt.Fprintf(w, body)
    fmt.Fprintf(w, footer)

}
