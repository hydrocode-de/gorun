.PHONY: frontend-build build run dev

frontend-build:
	cd manager && npm install && npm run build

build: frontend-build
	go build -o gorun .

run: build
	./gorun serve

dev:
	air serve