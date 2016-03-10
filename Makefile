all: build

build: n w g

n:
		go build -o ../nameserver-box github.com/toontong/box/cmd/nameserver/

w:
		go build -o ../worker-box github.com/toontong/box/cmd/worker/

g:
		go build -o ../gateway-box github.com/toontong/box/cmd/gateway/

install:
		go install github.com/toontong/box/cmd/nameserver/
		go install github.com/toontong/box/cmd/gateway/
		go install github.com/toontong/box/cmd/worker/

clean:
		go clean -i ./...github.com/toontong/box
