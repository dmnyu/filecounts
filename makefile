tidy:
	go mod tidy

build:
	go mod tidy
	go build -o files main/main.go

clean:
	rm files

test:
	go test

install:
	cp files /usr/local/bin
