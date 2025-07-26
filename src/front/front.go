package front

import (
	"net/http"
	"fmt"
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
