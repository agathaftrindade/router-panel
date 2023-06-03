.DEFAULT_GOAL := build

build:
	GOOS=linux GOARCH=amd64 go build -o bin/routerinfo-linux-386 main.go
	GOOS=linux GOARCH=arm64 go build -o bin/routerinfo-linux-arm64 main.go
	cp -r static/ bin/static

copy-to-rasp:
	rsync -r bin/ raspinha@192.168.15.8:/opt/router-info/ --delete

deploy: build copy-to-rasp
