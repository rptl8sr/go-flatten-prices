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