package front

import (
	"net/http"
	"fmt"
	"obliv/src/system"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
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

func HomePage(c *gin.Context) {
	html := `
	<body>
	<p>this is obliv homepage</p>
	</body>
	`
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


func CpuPage(c *gin.Context) {
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

	c.Writer.Write([]byte(fmt.Sprintf(html, system.PrintCPU())))
}

func RegisterPage(c *gin.Context) {

	if c.Request.Method == http.MethodPost {
		
		username := c.PostForm("username")
		password := c.PostForm("password")

		db := system.ConnectDatabase()
		defer db.Close()

		err := system.Register(db, username, password)
		if err != nil {
			c.Writer.Write([]byte(head))
			c.Writer.Write([]byte(fmt.Sprintf(`
			<body>
			<p style="color: red;">Failed: %s</p>
			<p><a href="/register">Try again</a></p>
			</body>`, err.Error())))
			c.Writer.Write([]byte(footer))
			return
		}
		c.Redirect(http.StatusSeeOther, "/login")
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

	c.Writer.Write([]byte(head))
	c.Writer.Write([]byte(body))
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

	r.GET("/register", RegisterPage)
	r.POST("/register", RegisterPage)

	r.GET("/login", LoginPage)
	r.POST("/login", LoginPage)

	r.GET("/logout", Logout)
}
