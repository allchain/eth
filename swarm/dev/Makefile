.PHONY: build cluster test

default: build

build:
	go build -o bin/swarm    github.com/galaxy/galaxy-eth/cmd/swarm
	go build -o bin/geth     github.com/galaxy/galaxy-eth/cmd/geth
	go build -o bin/bootnode github.com/galaxy/galaxy-eth/cmd/bootnode

cluster: build
	scripts/boot-cluster.sh

test:
	go test -v github.com/galaxy/galaxy-eth/swarm/...
