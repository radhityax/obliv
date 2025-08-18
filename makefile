src = src/main/main.go
out = delphinium

build:
	go build -o $(out) $(src)
run:
	go run $(src)
clean:
	rm $(out)
	rm -rf data/
