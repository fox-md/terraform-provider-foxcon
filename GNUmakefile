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
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

# test:
# 	go test -v -cover -timeout=120s -parallel=10 ./...

test: down pretest
	TF_ACC=1 go test ./... -count=1  -v -cover -timeout=120s -parallel=10

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

pretest: checkup populate

checkup: up
	@for i in {1..10}; do code=$$(curl -LI -u "admin:admin-secret" http://localhost:8081 -o /dev/null -w '%{http_code}' -s);if [ "$$code" -eq 200 ]; then echo "Schema is UP" && exit 0; else sleep 5; fi; done; echo "Schema is DOWN" && exit 1

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
	curl -s -o /dev/null -X PUT -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data '{"compatibility": "NONE"}'  http://localhost:8081/config/keep-latest ; \
	curl -s -o /dev/null -X PUT -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data '{"compatibility": "NONE"}'  http://localhost:8081/config/keep-active ; \
	curl -s -o /dev/null -X PUT -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data '{"compatibility": "NONE"}'  http://localhost:8081/config/switch-cleanup-mode ; \
	curl -s -o /dev/null -X PUT -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data '{"compatibility": "NONE"}'  http://localhost:8081/config/one ; \
	curl -s -o /dev/null -X PUT -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data '{"compatibility": "NONE"}'  http://localhost:8081/config/one-to-two-active ; \
	curl -s -o /dev/null -X PUT -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data '{"compatibility": "NONE"}'  http://localhost:8081/config/one-to-two-latest ; \
	curl -s -o /dev/null -X PUT -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data '{"compatibility": "NONE"}'  http://localhost:8081/config/data-source ; \


	schema=$$(cat tests/schemas/v1.json | tr -d '\n\r' | sed 's/"/\\"/g'); \
	payload="{\"schema\": \"$${schema}\", \"schemaType\": \"JSON\"}"; \
	curl -s -o /dev/null -X POST -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/one/versions ; \
	curl -s -o /dev/null -X POST -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/one-to-two-active/versions ; \
	curl -s -o /dev/null -X POST -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/one-to-two-latest/versions ; \

	for i in {1..5}; do \
		schema=$$(cat tests/schemas/v$$i.json | tr -d '\n\r' | sed 's/"/\\"/g'); \
		payload="{\"schema\": \"$${schema}\", \"schemaType\": \"JSON\"}"; \
		curl -s -o /dev/null -X POST -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/keep-latest/versions ; \
		sleep 0.5; \
		curl -s -o /dev/null -X POST -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/keep-active/versions ; \
		sleep 0.5; \
		curl -s -o /dev/null -X POST -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/switch-cleanup-mode/versions ; \
		sleep 0.5; \
		curl -s -o /dev/null -X POST -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/data-source/versions ; \
		sleep 0.5; \
	done ; \

	for i in {1..2}; do \
		curl -s -o /dev/null -X DELETE -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/keep-latest/versions/$$i ; \
		sleep 0.5; \
		curl -s -o /dev/null -X DELETE -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/switch-cleanup-mode/versions/$$i ; \
		sleep 0.5; \
		curl -s -o /dev/null -X DELETE -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/data-source/versions/$$i ; \
		sleep 0.5; \
	done ; \

	for i in {1..3}; do \
		curl -s -o /dev/null -X DELETE -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/keep-active/versions/$$i ; \
		sleep 0.5; \
	done ; \

add2oneactive:
	schema=$$(cat tests/schemas/v2.json | tr -d '\n\r' | sed 's/"/\\"/g'); \
	payload="{\"schema\": \"$${schema}\", \"schemaType\": \"JSON\"}"; \
	curl -s -o /dev/null -X POST -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/one-to-two-active/versions ; \

add2onelatest:
	schema=$$(cat tests/schemas/v2.json | tr -d '\n\r' | sed 's/"/\\"/g'); \
	payload="{\"schema\": \"$${schema}\", \"schemaType\": \"JSON\"}"; \
	curl -s -o /dev/null -X POST -u "admin:admin-secret" -H "Content-Type: application/vnd.schemaregistry.v1+json" --data "$$payload" http://localhost:8081/subjects/one-to-two-latest/versions ; \

clean:
	@rm -rf terraform.tfstate terraform.tfstate.*

docs:
	cd tools && go generate ./... && cd ..
	cd docs && \
	find . -type f -name "*.md" -exec sed -i 's/&#96;/`/g' {} + && \
	sed -i '/^description: |-/,/^---/c\description: |-\n\n---' index.md && \
	cd ..

.PHONY: fmt lint test testacc install generate restart up down apply plan docs

.SILENT: populate add2oneactive add2onelatest docs
