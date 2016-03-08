all: build install

build: n w

n:
		go build -o box_namserver github.com/toontong/box/cmd/nameserver/

w:
		go build -o box_worker github.com/toontong/box/cmd/worker/


install:
		go install github.com/toontong/box/cmd/nameserver/
		go install github.com/toontong/box/cmd/worker/

clean:
		go clean -i ./...github.com/toontong/box
