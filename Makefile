all: build

prepare:
	go-bindata-assetfs -pkg server -o server/assets.go web/...

build: prepare
	go build -o phs cmd/phs/main.go

clean:
	rm server/assets.go phs