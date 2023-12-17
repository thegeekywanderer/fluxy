help:
	@echo ''
	@echo 'Usage: make [TARGET] [EXTRA_ARGUMENTS]'
	@echo 'Targets:'
	@echo 'make dev: make dev for development work'
	@echo 'make build: make build container'
	@echo 'clean: clean for all clear docker images'

dev:
	docker-compose -f docker-compose.yml down
	if [ ! -f .env ]; then cp .env.example .env; fi;
	docker-compose -f docker-compose.yml up --build

build:
	docker build -f Dockerfile.build . 

clean:
	docker-compose -f docker-compose.yml down -v
