
targets: basic_subscriber basic_publisher ping pong

.PHONY: all basic_subscriber basic_publisher ping pong
all: $(targets)

ping: ping/ping
	cd $@
	go build

pong: pong/pong
	cd $@
	go build

basic_subscriber: basic_subscriber/basic_subscriber
	cd $@
	go build

basic_publisher: basic_publisher/basic_publisher
	cd $@
	go build
