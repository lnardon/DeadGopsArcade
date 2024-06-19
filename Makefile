.PHONY: all build

all: build

go.mod:
	go mod init jogo
	go get -u github.com/nsf/termbox-go

build: go.mod
	go build server.go map.go elemento.go utils.go interfaces.go clientState.go
	go build client.go map.go elemento.go utils.go interfaces.go clientState.go
	
clean:
	rm -f jogo

distclean: clean
	rm -f go.mod go.sum

build_client:
	go build ./client/*.go

build_server:
	go build ./server/*.go