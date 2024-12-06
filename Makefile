ifneq (,$(wildcard ./.env))
	include .env
	export $(shell sed 's/=.*//' .env)
endif

# Mock generation binary path
MOCKGEN := $(shell go env GOPATH)/bin/mockgen

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
## mock: generates mocks for all repositories
mock: mock_token_repository mock_user_repository # Add additional mocks here
	@echo "All mocks generated successfully!"

## mock_token_repository: generates mocks for the TokenRepository
mock_token_repository:
	@echo "Generating mock for TokenRepository..."
	$(MOCKGEN) -source=internal/pkg/repository/token_repository.go \
	           -destination=internal/mocks/mock_token_repository.go \
	           -package=mocks
	@echo "Mock for TokenRepository generated successfully!"

## mock_user_repository: generates mocks for the UserRepository
mock_user_repository:
	@echo "Generating mock for UserRepository..."
	$(MOCKGEN) -source=internal/pkg/repository/user_repository.go \
	           -destination=internal/mocks/mock_user_repository.go \
	           -package=mocks
	@echo "Mock for UserRepository generated successfully!"

# Add more mock generation rules as needed

# Variables for database connection
# Variables for soda command
SODA_CMD=/opt/homebrew/bin/soda  # Correct path to the soda binary

# Target to run migrations
migrate:
	$(SODA_CMD) migrate up

# Target to run the rollback (optional)
rollback:
	$(SODA_CMD) migrate down

# Target to create a new migration (optional)
new-migration:
	$(SODA_CMD) generate migration $(MIGRATION_NAME)

# Target to apply migrations
migrate-up: migrate

# Target to create the database (if needed)
create-db:
	psql -c "CREATE DATABASE usualstore;" -U postgres

# Target to drop the database (if needed)
drop-db:
	psql -c "DROP DATABASE usualstore;" -U postgres

.PHONY: migrate rollback new-migration migrate-up create-db drop-db