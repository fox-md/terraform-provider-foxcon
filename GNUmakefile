#default: fmt lint install generate

# build:
# 	go build -v ./...

# install: build
# 	go install -v ./...

install:
	go install .

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

# test:
# 	go test -v -cover -timeout=120s -parallel=10 ./...

test:
	TF_ACC=1 go test ./... -count=1  -v

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

restart: down up

up:
	docker-compose up -d

down:
	docker-compose stop
	docker-compose rm -f

apply:
	TF_LOG_PROVIDER=DEBUG terraform apply --auto-approve

plan:
	TF_LOG_PROVIDER=DEBUG terraform plan

.PHONY: fmt lint test testacc install generate restart up down apply plan
