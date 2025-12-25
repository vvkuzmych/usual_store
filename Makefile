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

## build_invoice: builds the invoice
build_invoice:
	@echo "Building invoice..."
	@go build -o dist/invoice ./cmd/micro/invoice
	@echo "Invoice built!"


## build_back: builds the back end
build_back:
	@echo "Building back end..."
	@go build -o dist/usualstore_api ./cmd/api
	@echo "Back end built!"

## start: starts front and back end
start: start_front start_back start_invoice

## start_front: starts the front end
start_front: build_front
	@echo "Starting the front end..."
	@env STRIPE_KEY=${STRIPE_KEY} STRIPE_SECRET=${STRIPE_SECRET} ./dist/usualstore -port=${USUAL_STORE_PORT} &
	@echo "Front end running!"

## start_invoice: starts invoice microservice
start_invoice: build_invoice
	@echo "Starting the invoice.."
	@./dist/invoice &
	@echo "invoice running!"

## start_back: starts the back end
start_back: build_back
	@echo "Starting the back end..."
	@env STRIPE_KEY=${STRIPE_KEY} STRIPE_SECRET=${STRIPE_SECRET} ./dist/usualstore_api -port=${API_PORT} &
	@echo "Back end running!"

## stop: stops the front and back end
stop: stop_front stop_back stop_invoice
	@echo "All applications stopped"

## stop_front: stops the front end
stop_front:
	@echo "Stopping the front end..."
	@-pkill -SIGTERM -f "usualstore -port=${USUAL_STORE_PORT}"
	@echo "Stopped front end"

## stop_invoice: stops invoice microservice
stop_invoice:
	@echo "Stopping invoice microservice..."
	@-pkill -SIGTERM -f "invoice"
	@echo "Stopped invoice microservice"


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

## Docker Compose targets
docker-up:
	@echo "Starting all services with Docker Compose..."
	docker compose up -d
	@echo "Services started!"

docker-down:
	@echo "Stopping all services..."
	docker compose down
	@echo "Services stopped!"

docker-restart:
	@echo "Restarting all services..."
	docker compose restart
	@echo "Services restarted!"

docker-logs:
	@echo "Following logs..."
	docker compose logs -f

docker-ps:
	@echo "Listing running services..."
	docker compose ps

## Frontend application targets
react-start:
	@echo "🚀 Starting React Frontend + Backend..."
	@docker compose --profile react-frontend up -d
	@echo ""
	@echo "✅ React Frontend running at: http://localhost:3000"
	@echo "✅ Backend API running at: http://localhost:4001"

typescript-start:
	@echo "🚀 Starting TypeScript Frontend + Backend..."
	@docker compose --profile typescript-frontend up -d
	@echo ""
	@echo "✅ TypeScript Frontend running at: http://localhost:3001"
	@echo "✅ Backend API running at: http://localhost:4001"

go-start:
	@echo "🚀 Starting Go Frontend + Backend..."
	@docker compose --profile go-frontend up -d
	@echo ""
	@echo "✅ Go Frontend running at: http://localhost:4000"
	@echo "✅ Backend API running at: http://localhost:4001"

all-frontends-start:
	@echo "🚀 Starting ALL Frontends + Backend..."
	@docker compose --profile go-frontend --profile react-frontend --profile typescript-frontend up -d
	@echo ""
	@echo "✅ Go Frontend running at:         http://localhost:4000"
	@echo "✅ React Frontend running at:      http://localhost:3000"
	@echo "✅ TypeScript Frontend running at: http://localhost:3001"
	@echo "✅ Backend API running at:         http://localhost:4001"

react-stop:
	@echo "🛑 Stopping React Frontend..."
	@docker compose --profile react-frontend down
	@echo "✅ React Frontend stopped!"

typescript-stop:
	@echo "🛑 Stopping TypeScript Frontend..."
	@docker compose --profile typescript-frontend down
	@echo "✅ TypeScript Frontend stopped!"

go-stop:
	@echo "🛑 Stopping Go Frontend..."
	@docker compose --profile go-frontend down
	@echo "✅ Go Frontend stopped!"

react-logs:
	@echo "📋 React Frontend logs:"
	@docker compose logs -f react-frontend

typescript-logs:
	@echo "📋 TypeScript Frontend logs:"
	@docker compose logs -f typescript-frontend

go-logs:
	@echo "📋 Go Frontend logs:"
	@docker compose logs -f go-frontend

react-build:
	@echo "🔨 Building React Frontend Docker image..."
	@docker compose --profile react-frontend build react-frontend
	@echo "✅ React Frontend image built!"

typescript-build:
	@echo "🔨 Building TypeScript Frontend Docker image..."
	@docker compose --profile typescript-frontend build typescript-frontend
	@echo "✅ TypeScript Frontend image built!"

go-build-docker:
	@echo "🔨 Building Go Frontend Docker image..."
	@docker compose --profile go-frontend build go-frontend
	@echo "✅ Go Frontend image built!"

build-all-frontends:
	@echo "🔨 Building all frontend Docker images..."
	@docker compose --profile go-frontend --profile react-frontend --profile typescript-frontend build
	@echo "✅ All frontend images built!"

react-restart:
	@echo "🔄 Restarting React Frontend..."
	@docker compose restart react-frontend
	@echo "✅ React Frontend restarted!"

typescript-restart:
	@echo "🔄 Restarting TypeScript Frontend..."
	@docker compose restart typescript-frontend
	@echo "✅ TypeScript Frontend restarted!"

go-restart:
	@echo "🔄 Restarting Go Frontend..."
	@docker compose restart go-frontend
	@echo "✅ Go Frontend restarted!"

frontend-status:
	@echo "📊 Frontend Services Status:"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@docker compose ps | grep -E "frontend|PORT" || echo "No frontends running"

## IPv6 connectivity targets
test-ipv6:
	@echo "Running IPv6 connectivity tests..."
	@./scripts/test-ipv6.sh

check-ipv6-network:
	@echo "Checking IPv6 network configuration..."
	@docker network inspect $$(docker compose config | grep -A 5 "networks:" | grep -v "networks:" | head -1 | tr -d ' ' | cut -d: -f1 | xargs -I {} echo "usual_store_{}" | sed 's/usual_store_/usualstore_/') | grep -A 10 IPv6 || echo "IPv6 not configured"

verify-db-ipv6:
	@echo "Verifying PostgreSQL IPv6 connectivity..."
	@docker compose exec database ss -tln | grep -E '::1|:::' | grep 5432 && echo "✅ PostgreSQL listening on IPv6" || echo "❌ PostgreSQL NOT listening on IPv6"

test-db-ipv6-host:
	@echo "Testing database connection from host via IPv6..."
	@psql "postgres://postgres:password@[::1]:5432/usualstore?sslmode=disable" -c "SELECT '✅ IPv6 connection successful!' as status;" 2>/dev/null || echo "❌ IPv6 connection failed"

test-db-ipv4-host:
	@echo "Testing database connection from host via IPv4..."
	@psql "postgres://postgres:password@127.0.0.1:5432/usualstore?sslmode=disable" -c "SELECT '✅ IPv4 connection successful!' as status;" 2>/dev/null || echo "❌ IPv4 connection failed"

show-container-ips:
	@echo "Container IPv4 and IPv6 addresses:"
	@echo "===================================="
	@for service in database back-end front-end invoice; do \
		container=$$(docker compose ps -q $$service); \
		if [ -n "$$container" ]; then \
			echo "\n$$service:"; \
			docker inspect $$container | grep -A 5 Networks | grep -E 'IPAddress|GlobalIPv6Address' | sed 's/^/  /'; \
		fi; \
	done

## Database helpers with IPv6 support
db-shell-ipv6:
	@echo "Connecting to PostgreSQL via IPv6..."
	@psql "postgres://postgres:password@[::1]:5432/usualstore?sslmode=disable"

db-shell-ipv4:
	@echo "Connecting to PostgreSQL via IPv4..."
	@psql "postgres://postgres:password@127.0.0.1:5432/usualstore?sslmode=disable"

db-shell-docker:
	@echo "Connecting to PostgreSQL inside Docker..."
	@docker compose exec database psql -U postgres -d usualstore

## Help target
help:
	@echo "Available targets:"
	@echo ""
	@echo "Build targets:"
	@echo "  build                  - Build all binaries"
	@echo "  build_front            - Build front-end binary"
	@echo "  build_back             - Build back-end binary"
	@echo "  build_invoice          - Build invoice binary"
	@echo "  clean                  - Clean all binaries"
	@echo ""
	@echo "Run targets:"
	@echo "  start                  - Start all services"
	@echo "  stop                   - Stop all services"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-up              - Start services with Docker Compose"
	@echo "  docker-down            - Stop Docker Compose services"
	@echo "  docker-restart         - Restart Docker Compose services"
	@echo "  docker-logs            - Follow Docker Compose logs"
	@echo "  docker-ps              - List running containers"
	@echo ""
	@echo "Frontend targets:"
	@echo "  react-start            - Start React Frontend (port 3000) + Backend"
	@echo "  typescript-start       - Start TypeScript Frontend (port 3001) + Backend"
	@echo "  go-start               - Start Go Frontend (port 4000) + Backend"
	@echo "  all-frontends-start    - Start ALL frontends + Backend"
	@echo "  react-stop             - Stop React Frontend"
	@echo "  typescript-stop        - Stop TypeScript Frontend"
	@echo "  go-stop                - Stop Go Frontend"
	@echo "  react-logs             - View React Frontend logs"
	@echo "  typescript-logs        - View TypeScript Frontend logs"
	@echo "  go-logs                - View Go Frontend logs"
	@echo "  react-build            - Build React Frontend Docker image"
	@echo "  typescript-build       - Build TypeScript Frontend Docker image"
	@echo "  go-build-docker        - Build Go Frontend Docker image"
	@echo "  build-all-frontends    - Build all frontend Docker images"
	@echo "  react-restart          - Restart React Frontend"
	@echo "  typescript-restart     - Restart TypeScript Frontend"
	@echo "  go-restart             - Restart Go Frontend"
	@echo "  frontend-status        - Show status of all frontends"
	@echo ""
	@echo "IPv6 targets:"
	@echo "  test-ipv6              - Run comprehensive IPv6 connectivity tests"
	@echo "  check-ipv6-network     - Check Docker network IPv6 configuration"
	@echo "  verify-db-ipv6         - Verify PostgreSQL IPv6 listening"
	@echo "  test-db-ipv6-host      - Test PostgreSQL connection via IPv6 from host"
	@echo "  test-db-ipv4-host      - Test PostgreSQL connection via IPv4 from host"
	@echo "  show-container-ips     - Show all container IPv4/IPv6 addresses"
	@echo ""
	@echo "Database targets:"
	@echo "  db-shell-ipv6          - Connect to PostgreSQL via IPv6"
	@echo "  db-shell-ipv4          - Connect to PostgreSQL via IPv4"
	@echo "  db-shell-docker        - Connect to PostgreSQL inside Docker"
	@echo "  migrate                - Run database migrations"
	@echo "  rollback               - Rollback database migrations"
	@echo "  create-db              - Create database"
	@echo "  drop-db                - Drop database"
	@echo ""
	@echo "Other targets:"
	@echo "  mock                   - Generate all mocks"
	@echo "  help                   - Show this help message"

.PHONY: migrate rollback new-migration migrate-up create-db drop-db
.PHONY: docker-up docker-down docker-restart docker-logs docker-ps
.PHONY: react-start typescript-start go-start all-frontends-start
.PHONY: react-stop typescript-stop go-stop
.PHONY: react-logs typescript-logs go-logs
.PHONY: react-build typescript-build go-build-docker build-all-frontends
.PHONY: react-restart typescript-restart go-restart
.PHONY: frontend-status
.PHONY: test-ipv6 check-ipv6-network verify-db-ipv6 test-db-ipv6-host test-db-ipv4-host show-container-ips
.PHONY: db-shell-ipv6 db-shell-ipv4 db-shell-docker help
