SHELL := /bin/bash

#default: fmt lint install generate

build:
	go build -v ./...

# install: build
# 	go install -v ./...

install:
	go install .

lint:
	act -j lint
#golangci-lint run

generate:
	cd tools; go generate ./...; cp -r ../examples/guides ../docs

fmt:
	gofmt -s -w -e .

# test:
# 	go test -v -cover -timeout=120s -parallel=10 ./...

test:
	TF_ACC=1 go test ./... -count=1 -run='$(TEST)' -v

tests: pretest
	TF_ACC=1 go test ./... -count=1  -v -cover -timeout=120s -parallel=10

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

pretest: restart populate

restart: down up
	@for i in {1..25}; do code=$$(curl -LI -u "admin:admin-secret" http://localhost:8081 -o /dev/null -w '%{http_code}' -s);if [ "$$code" -eq 200 ]; then echo "Schema is UP" && exit 0; else sleep 5; fi; done; echo "Schema is DOWN" && exit 1

up:
	docker-compose up -d

down:
	docker-compose stop
	docker-compose rm -f

apply:
	TF_LOG_PROVIDER=DEBUG terraform apply --auto-approve

plan:
	TF_LOG_PROVIDER=DEBUG terraform plan

populate:
	echo "Preparing testdata"
	curl -s -o /dev/null -X PUT -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data '{"compatibility": "NONE"}'  http://localhost:8081/config/one ; \

	schema=$$(cat tests/schemas/v1.json | tr -d '\n\r' | sed 's/"/\\"/g'); \
	payload="{\"schema\": \"$${schema}\", \"schemaType\": \"JSON\"}"; \
	curl -s -o /dev/null -X POST -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/one/versions ; \

clean:
	@rm -rf terraform.tfstate terraform.tfstate.*

docs:
	cd tools && go generate ./... && cd ..
	cd docs && \
	find . -type f -name "*.md" -exec sed -i 's/&#96;/`/g' {} + && \
	cd ..
	cp -r examples/guides docs/

.PHONY: fmt lint test testacc install generate restart up down apply plan docs

.SILENT: populate docs
