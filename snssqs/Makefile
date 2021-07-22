.PHONY: build pub sub

default: build

build: pub sub

pub:
	go build -o pub pub.go helpers.go

sub:
	go build -o sub sub.go helpers.go
