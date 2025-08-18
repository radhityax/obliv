package front

import (
	"net/http"
	"fmt"
	"delphinium/src/system"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)

var head = `
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>delphinium</title>
</head>
`

var footer = `
<footer><p>powered by delphinium - open source linux control panel</p></footer>
</html>
`

var back = "<a href=\"/\">back</a>"

func GetUsername(c *gin.Context) string {
	session := sessions.Default(c)
	return session.Get("user").(string)
}
func HomePage(c *gin.Context) {
	username := GetUsername(c)

	html := fmt.Sprintf(`
	<body>
	<p>this is delphinium homepage</p>
	<p>hello, %s. (<a href="/logout">logout</a>)</p>
	<ul>
	<li><a href="/memory">memory</a></li>
	<li><a href="/cpu">cpu</a></li>
	</ul>
	</body>
	`, username)
	c.Writer.Write([]byte(head))
	c.Writer.Write([]byte(html))
	c.Writer.Write([]byte(footer))
}

func MemoryPage(c *gin.Context) {
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
	memUsed := dat.MemTotal - dat.MemFree - dat.Buffers - dat.Cached

	c.Writer.Write([]byte(head))
	c.Writer.Write([]byte(fmt.Sprintf(html, dat.MemTotal, memUsed,
	dat.MemAvailable, dat.Buffers, dat.Cached)))
	c.Writer.Write([]byte(back))
	c.Writer.Write([]byte(footer))
}

func MemoryData(c *gin.Context) {
	dat := system.GetMemory()
	response := map[string]interface{}{
		"memTotal":    dat.MemTotal,
		"memUsed":     dat.MemTotal - dat.MemFree - dat.Buffers - dat.Cached,
		"memAvailable": dat.MemAvailable,
		"buffers":     dat.Buffers,
		"cached":      dat.Cached,
	}
	c.JSON(http.StatusOK, response)
}

func CpuData(c *gin.Context) {
	stats, _ := system.GetCpuStats()
	loadone, loadtwo, loadthree, _ := system.GetLoadAverage()

	response := map[string]interface{}{
		"User": stats.User,
		"Nice": stats.Nice,
		"System": stats.System,
		"Idle": stats.Idle,
		"Iowait": stats.Iowait,
		"Irq": stats.Irq,
		"Softirq": stats.Softirq,
		"Steal": stats.Steal,
		"Guest": stats.Guest,
		"GuestNice": stats.GuestNice,
		"LoadOne": loadone,
		"LoadTwo": loadtwo,
		"LoadThree": loadthree,
	}
	c.JSON(http.StatusOK, response)
}


func CpuPage(c *gin.Context) {
	data, _ := system.GetCpuStats()
	loadone, loadtwo, loadthree, _ := system.GetLoadAverage()
	html := `
	<script>
	function updatecpu() {
		fetch('/cpu-data')
		.then(response => response.json())
		.then(data => {
			document.getElementById("User").innerText = 'User: ' + data.User;
			document.getElementById("Nice").innerText = 'Nice: ' + data.Nice;
			document.getElementById("System").innerText = 'System: ' + data.System;
			document.getElementById("Idle").innerText = 'Idle: ' + data.Idle;
			document.getElementById("Iowait").innerText = 'Iowait: ' + data.Iowait;
			document.getElementById("Irq").innerText = 'Irq: ' + data.Irq;
			document.getElementById("Softirq").innerText = 'Softirq: ' + data.Softirq;
			document.getElementById("Steal").innerText = 'Steal: ' + data.Steal;
			document.getElementById("Guest").innerText = 'Guest: ' + data.Guest;
			document.getElementById("GuestNice").innerText = 'GuestNice: ' + data.GuestNice;
			document.getElementById("LoadOne").innerText = 'Load One: ' + data.LoadOne;
			document.getElementById("LoadTwo").innerText = 'Load Two: ' + data.LoadTwo;
			document.getElementById("LoadThree").innerText = 'Load Three: ' + data.LoadThree;
		});
	}
	setInterval(updatecpu, 3000);
	</script>
	</head>
	<body>
	<p>this is cpu</p>
	<p id="User">User: %d</p>
	<p id="Nice">Nice: %d</p>
	<p id="System">System: %d</p>
	<p id="Idle">Idle: %d</p>
	<p id="Iowait">Iowait: %d</p>
	<p id="Irq">Irq: %d</p>
	<p id="Softirq">Softirq: %d</p>
	<p id="Steal">Steal: %d</p>
	<p id="Guest">Guest: %d</p>
	<p id="GuestNice">GuestNice: %d</p>
	<p id="LoadOne">Load One: %f</p>
	<p id="LoadTwo">Load Two: %f</p>
	<p id="LoadThree">Load Three: %f</p>
	</body>
	</html>
	`

	c.Writer.Write([]byte(head))
	c.Writer.Write([]byte(fmt.Sprintf(html, data.User, data.Nice, data.System,
	data.Idle, data.Iowait, data.Irq, data.Softirq, data.Steal, data.Guest,
	data.GuestNice, loadone, loadtwo, loadthree)))
	c.Writer.Write([]byte(back))
	c.Writer.Write([]byte(footer))
}
func LoginPage(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		username := c.PostForm("username")
		password := c.PostForm("password")

		db := system.ConnectDatabase()
		defer db.Close()

		var hashedPassword string
		query := `SELECT password FROM account WHERE username = ?`
		err := db.QueryRow(query, username).Scan(&hashedPassword)
		if err != nil {
			c.String(http.StatusUnauthorized, "Username ga ketemu")
			return
		}

		if !system.CheckPassword(hashedPassword, password) {
			c.String(http.StatusUnauthorized, "Invalid passwordnya")
			return
        }

		session := sessions.Default(c)
		session.Set("user", username)
		session.Save()

		c.Redirect(http.StatusSeeOther, "/")
		return
    }

	loginForm := `
	<form method="POST">
	<input name="username" placeholder="Username" required>
	<input name="password" type="password" placeholder="Password" required>
	<button type="submit">Login</button>
	</form>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, loginForm)
	c.Writer.Write([]byte(footer))
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusSeeOther, "/login")
}

func FrontSetup(r *gin.Engine) {

	r.GET("/", system.AuthRequired(), HomePage)
	r.GET("/memory", system.AuthRequired(), MemoryPage)
	r.GET("/memory-data",  system.AuthRequired(), MemoryData)
	r.GET("/cpu", system.AuthRequired(),CpuPage)
	r.GET("/cpu-data", system.AuthRequired(),CpuData)
	r.GET("/login", LoginPage)
	r.POST("/login", LoginPage)

	r.GET("/logout", Logout)
}
