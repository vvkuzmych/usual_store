ifneq (,$(wildcard ./.env))
	include .env
	export $(shell sed 's/=.*//' .env)
endif

run:
	echo "Starting the application..."
	@echo "Using Stripe Secret: $(STRIPE_SECRET)"
	@echo "Using Stripe Key: $(STRIPE_KEY)"
	@echo "Using Usual Store Port: $(USUAL_STORE_PORT)"
	@echo "Using API Port: $(API_PORT)"

## build: builds all binaries
build: clean build_front build_back
	@printf "All binaries built!\n"

## clean: cleans all binaries and runs go clean
clean:
	@echo "Cleaning..."
	@- rm -f dist/*
	@go clean
	@echo "Cleaned!"

## build_front: builds the front end
build_front:
	@echo "Building front end..."
	@go build -o dist/usualstore ./cmd/web
	@echo "Front end built!"

## build_back: builds the back end
build_back:
	@echo "Building back end..."
	@go build -o dist/usualstore_api ./cmd/api
	@echo "Back end built!"

## start: starts front and back end
start: start_front start_back

## start_front: starts the front end
start_front: build_front
	@echo "Starting the front end..."
	@env STRIPE_KEY=${STRIPE_KEY} STRIPE_SECRET=${STRIPE_SECRET} ./dist/usualstore -port=${USUAL_STORE_PORT} &
	@echo "Front end running!"

## start_back: starts the back end
start_back: build_back
	@echo "Starting the back end..."
	@env STRIPE_KEY=${STRIPE_KEY} STRIPE_SECRET=${STRIPE_SECRET} ./dist/usualstore_api -port=${API_PORT} &
	@echo "Back end running!"

## stop: stops the front and back end
stop: stop_front stop_back
	@echo "All applications stopped"

## stop_front: stops the front end
stop_front:
	@echo "Stopping the front end..."
	@-pkill -SIGTERM -f "usualstore -port=${USUAL_STORE_PORT}"
	@echo "Stopped front end"

## stop_back: stops the back end
stop_back:
	@echo "Stopping the back end..."
	@-pkill -SIGTERM -f "usualstore_api -port=${API_PORT}"
	@echo "Stopped back end"