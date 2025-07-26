src = src/main.go
out = obliv

build:
	go build -o $(out) $(src)
run:
	go run $(src)
