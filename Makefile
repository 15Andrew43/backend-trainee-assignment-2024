.PHONY: all build run stop

all: build run

build:
	docker-compose -f ./docker-compose.yml build

run:
	docker-compose -f ./docker-compose.yml up -d

stop:
	docker-compose -f ./docker-compose.yml down