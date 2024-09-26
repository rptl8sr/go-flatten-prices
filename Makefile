PROJECT := "go-flattern-prices"
V := "0.0.1"
USER := rptl8sr
EMAIL := $(USER)@gmail.com
LOCAL_BIN :=$(CURDIR)/bin
MIGRATIONS_DIR := migrations


.PHONY: git-init
git-init:
	gh repo create $(PROJECT) --private
	git init
	git config user.name "$(USER)"
	git config user.email "$(EMAIL)"
	git add Makefile go.mod
	git commit -m "Init commit"
	git remote add origin git@github.com:$(USER)/$(PROJECT).git
	git remote -v
	git push -u origin master


BN ?= dev
# make git-checkout BN=dev
.PHONY: git-checkout
git-checkout:
	git checkout -b $(BN)


.PHONY: blueprint
blueprint:
	mkdir -p cmd && echo 'package main\n\nfunc main() {\n}\n' > cmd/main.go
	mkdir -p app && touch app/app.go
	mkdir -p input
	mkdir -p output
	mkdir -p logs
	mkdir -p config && touch config/config.yaml
	mkdir -p internal/logger && touch internal/logger/logger.go
	mkdir -p internal/store && touch internal/store/store.go
	mkdir -p internal/controller && touch internal/controller/controller.go
	mkdir -p internal/processors && touch internal/processors/processors.go && touch internal/processors/images.go && touch internal/processors/prices.go


.PHONY: lint
lint:
	golangci-lint run ./...


.PHONY: test
test:
	go vet ./...
	go test ./...


.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(LOCAL_BIN)/$(PROJECT)_$(V) ./cmd


.PHONY: build-win64
build-linux:
	GOOS=windows GOARCH=amd64 go build -o $(LOCAL_BIN)/$(PROJECT)_v$(V) ./cmd