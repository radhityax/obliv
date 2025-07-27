package front

import (
	"net/http"
	"fmt"
	
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
	<meta http-equiv="refresh" content="3" />
	<title>obliv</title>
	</head>
	<body>
	<p>this is memory</p>
	<p>total: %d</p>
	<p>Used: %d</p>
	<p>Available: %d</p>
	<p>Buff/Cache: %d/%d</p>
	</body>
	</html>
	`

	fmt.Fprintf(w, html, dat.MemTotal, dat.MemTotal - dat.MemFree - 
	dat.Buffers - dat.Cached, dat.MemAvailable, dat.Buffers, dat.Cached)
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
