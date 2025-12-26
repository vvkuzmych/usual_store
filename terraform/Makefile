.PHONY: help init plan apply destroy clean validate fmt lint test policy-test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

init: ## Initialize Terraform
	terraform init

plan: ## Show what Terraform will do
	terraform plan

apply: ## Apply Terraform configuration
	terraform apply

apply-auto: ## Apply without confirmation (use with caution!)
	terraform apply -auto-approve

destroy: ## Destroy all infrastructure
	terraform destroy

destroy-auto: ## Destroy without confirmation (use with caution!)
	terraform destroy -auto-approve

clean: ## Clean Terraform files
	rm -rf .terraform .terraform.lock.hcl
	rm -f terraform.tfstate terraform.tfstate.backup
	rm -f crash.log

validate: ## Validate Terraform configuration
	terraform validate

fmt: ## Format Terraform files
	terraform fmt -recursive

lint: ## Lint Terraform files (requires tflint)
	@command -v tflint >/dev/null 2>&1 || { echo "tflint not installed. Run: brew install tflint"; exit 1; }
	tflint

test: ## Run Terraform tests
	terraform test

policy-test: ## Test OPA policies
	@command -v opa >/dev/null 2>&1 || { echo "OPA not installed. Run: brew install opa"; exit 1; }
	cd modules/policies/policies && opa test . -v

policy-fmt: ## Format OPA policy files
	@command -v opa >/dev/null 2>&1 || { echo "OPA not installed. Run: brew install opa"; exit 1; }
	cd modules/policies/policies && opa fmt -w *.rego

policy-check: ## Check OPA policies for errors
	@command -v opa >/dev/null 2>&1 || { echo "OPA not installed. Run: brew install opa"; exit 1; }
	cd modules/policies/policies && opa check *.rego

show: ## Show current Terraform state
	terraform show

output: ## Show all outputs
	terraform output

status: ## Show service status
	@echo "=== Service Status ==="
	@terraform output -json service_status 2>/dev/null | jq -r 'to_entries[] | "\(.key): \(.value)"' || echo "Run 'terraform apply' first"

urls: ## Show all service URLs
	@echo "=== Frontend URLs ==="
	@terraform output -json frontend_urls 2>/dev/null | jq -r 'to_entries[] | "\(.key): \(.value)"' || echo "Run 'terraform apply' first"
	@echo ""
	@echo "=== Backend URLs ==="
	@echo "API: $$(terraform output -raw backend_api_url 2>/dev/null || echo 'N/A')"
	@echo "Support: $$(terraform output -raw support_service_url 2>/dev/null || echo 'N/A')"
	@echo "Kafka UI: $$(terraform output -raw kafka_ui_url 2>/dev/null || echo 'N/A')"
	@echo "Jaeger: $$(terraform output -raw jaeger_ui_url 2>/dev/null || echo 'N/A')"
	@echo "OPA: $$(terraform output -raw policy_server_url 2>/dev/null || echo 'N/A')"

logs-database: ## Show database logs
	docker logs usualstore-database -f

logs-kafka: ## Show Kafka logs
	docker logs usualstore-kafka -f

logs-opa: ## Show OPA server logs
	docker logs usualstore-opa-server -f

logs-enforcer: ## Show policy enforcer logs
	docker logs usualstore-policy-enforcer -f

health-check: ## Check health of all services
	@echo "=== Health Check ==="
	@echo -n "Database: " && docker inspect usualstore-database --format='{{.State.Health.Status}}' 2>/dev/null || echo "not running"
	@echo -n "Kafka: " && docker inspect usualstore-kafka --format='{{.State.Health.Status}}' 2>/dev/null || echo "not running"
	@echo -n "Zookeeper: " && docker inspect usualstore-zookeeper --format='{{.State.Health.Status}}' 2>/dev/null || echo "not running"
	@echo -n "Jaeger: " && docker inspect usualstore-jaeger --format='{{.State.Health.Status}}' 2>/dev/null || echo "not running"
	@echo -n "OPA Server: " && docker inspect usualstore-opa-server --format='{{.State.Health.Status}}' 2>/dev/null || echo "not running"

policy-query: ## Query OPA policy (requires POLICY and INPUT env vars)
	@if [ -z "$(POLICY)" ]; then echo "Usage: make policy-query POLICY=usualstore/network/allow INPUT='{...}'"; exit 1; fi
	@curl -s -X POST http://localhost:8181/v1/data/$(POLICY) \
		-H "Content-Type: application/json" \
		-d '{"input": $(INPUT)}' | jq

network-check: ## Check if service communication is allowed
	@echo "Usage: make network-check SOURCE=react-frontend TARGET=back-end"
	@if [ -z "$(SOURCE)" ] || [ -z "$(TARGET)" ]; then \
		echo "Example: make network-check SOURCE=react-frontend TARGET=back-end"; \
		exit 1; \
	fi
	@curl -s -X POST http://localhost:8181/v1/data/usualstore/network/allow \
		-H "Content-Type: application/json" \
		-d '{"input": {"source_service": "$(SOURCE)", "target_service": "$(TARGET)"}}' | jq

security-scan: ## Run security scan on all containers
	@echo "=== Security Scan ==="
	@for container in $$(docker ps --format '{{.Names}}' | grep usualstore); do \
		echo "Scanning $$container..."; \
		curl -s -X POST http://localhost:8080/enforce -d "{\"container_name\": \"$$container\"}" | jq; \
	done

backup-state: ## Backup Terraform state
	@mkdir -p backups
	@cp terraform.tfstate backups/terraform.tfstate.$$(date +%Y%m%d_%H%M%S)
	@echo "State backed up to backups/"

restore-state: ## Restore Terraform state (requires BACKUP env var)
	@if [ -z "$(BACKUP)" ]; then echo "Usage: make restore-state BACKUP=terraform.tfstate.20250101_120000"; exit 1; fi
	@cp backups/$(BACKUP) terraform.tfstate
	@echo "State restored from $(BACKUP)"

graph: ## Generate Terraform dependency graph
	terraform graph | dot -Tpng > graph.png
	@echo "Graph saved to graph.png"

cost-estimate: ## Estimate infrastructure costs (placeholder for local)
	@echo "This is a local Docker setup - no cloud costs!"
	@echo "Resource usage:"
	@docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}" | grep usualstore || echo "No containers running"

install-tools: ## Install required tools (macOS)
	@echo "Installing required tools..."
	@command -v terraform >/dev/null 2>&1 || brew install terraform
	@command -v opa >/dev/null 2>&1 || brew install opa
	@command -v jq >/dev/null 2>&1 || brew install jq
	@command -v tflint >/dev/null 2>&1 || brew install tflint
	@echo "All tools installed!"

setup: install-tools init ## Complete setup (install tools + init)

quick-start: setup apply urls ## Quick start - setup and deploy everything

rebuild-enforcer: ## Rebuild policy enforcer
	cd modules/policies/policy-enforcer && docker build -t usualstore/policy-enforcer:latest .

update-policies: ## Update OPA policies without restarting
	@echo "Reloading policies..."
	@docker exec usualstore-opa-server opa build modules/policies/policies/
	@echo "Policies updated!"

compliance-report: ## Generate compliance report
	@echo "=== Compliance Report ===" > compliance-report.txt
	@echo "Generated: $$(date)" >> compliance-report.txt
	@echo "" >> compliance-report.txt
	@curl -s http://localhost:8080/audit >> compliance-report.txt
	@echo "Report saved to compliance-report.txt"

