ENV_GITLABCI_FILE := .env.gitlabci
ENV_GITLABCI = $(shell cat $(ENV_GITLABCI_FILE))
# コンマ区切りで go の package を列挙する
COVERAGE_TARGET_PACKAGES = gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database/model,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/adapter/database/repository,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/handler,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/router,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/api/schema,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/config,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/constructor,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entconst,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/entity,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/domain/validation,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/log,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/input_port,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/interactor,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port,gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/utils/set

# App Server
.PHONY: run-dev
run-dev:
	docker compose -f docker-compose.dev.yml up --build

.PHONY: destroy
destroy:
	docker compose -f docker-compose.dev.yml -f docker-compose.test.yml down --volumes --remove-orphans
	rm -rf mysql/mysql-data

# MySQL
.PHONY: init-db
init-db:
	docker compose -f docker-compose.dev.yml run --rm api go run cmd/initdb/initdb.go

# Tools
.PHONY: tools
tools:
	go install github.com/cosmtrek/air@v1.49.0 # go 1.21以前に対応するバージョン
	go install github.com/golang/mock/mockgen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1 # go 1.19に対応するバージョン

# Lint, Format
.PHONY: lint
lint: tools
	golangci-lint run ./... --timeout=5m

.PHONY: format
format: tools
	golangci-lint run ./... --fix

.PHONY: test
test:
	docker compose -f docker-compose.test.yml run --rm test-api go test -v ./...

.PHONY: test-without-db
test-without-db:
	docker compose -f docker-compose.test.yml run --rm test-api go test `go list ./... | grep -v gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/integration_tests`

.PHONY: test-gitlab
test-gitlab:
	go test -v -cover -coverpkg $(COVERAGE_TARGET_PACKAGES) ./... -coverprofile=cover.out && go tool cover -func cover.out

.PHONY: test-coverage
test-coverage:
	docker compose -f docker-compose.test.yml run --rm test-api  sh -c 'go test -v -cover -coverpkg $(COVERAGE_TARGET_PACKAGES) ./... -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html && go tool cover -func cover.out'

.PHONY: migrate
migrate:
	docker compose -f docker-compose.dev.yml run --rm api sh -c 'migrate -path mysql/migrations -database "mysql://$$DB_USER:$$DB_PASSWORD@tcp($$DB_HOST:$$DB_PORT)/$$DB_NAME?multiStatements=true" up'

.PHONY: check
check:
	echo "called"
