.PHONY: frontend-build build run dev

frontend-build:
	cd manager && npm install && npm run build

build: frontend-build
	go build -ldflags "-X github.com/hydrocode-de/gorun/version.Commit=$(shell git rev-parse HEAD) -X github.com/hydrocode-de/gorun/version.Date=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" -o gorun .

run: build
	./gorun serve

dev:
	air serve