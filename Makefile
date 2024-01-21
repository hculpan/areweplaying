all: build

build:
	templ generate
	go build -o areweplaying cmd/web/*.go
	go build -o awp-cli cmd/cli/*.go

linuxbuild:
	templ generate
	GOOS=linux GOARCH=amd64 go build -o areweplaying.linux cmd/web/*.go
	GOOS=linux GOARCH=amd64 go build -o awp-cli.linux cmd/cli/*.go

clean:
	rm -rf areweplaying
	rm -rf areweplaying.*

test:
