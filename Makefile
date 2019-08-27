HAS_GOLINT := $(shell command -v golint;)
VERSION :=
COMMIT :=
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
LDFLAGS := "-X main.version=${VERSION} -X main.commit=${COMMIT}"
BINARY := dockerctl
MAIN := ./cmd/dockerctl

.PHONY: install
install: bootstrap test build

.PHONY: test
test: golint
	go test ./... -v

.PHONY: gofmt
gofmt:
	gofmt -s -w .

.PHONY: golint
golint: gofmt
ifndef HAS_GOLINT
	go get -u golang.org/x/lint/golint
endif
	golint -set_exit_status ./cmd/...
	golint -set_exit_status ./pkg/...

.PHONY: build
build: clean bootstrap
	mkdir -p $(BUILD)
	go build -o $(BUILD)/$(BINARY) $(MAIN)

.PHONY: dist
dist:
ifeq ($(strip $(VERSION)),)
	$(error VERSION is not set)
endif
ifeq ($(strip $(COMMIT)),)
	$(error COMMIT is not set)
endif
	go get -u github.com/inconshreveable/mousetrap
	mkdir -p $(BUILD)
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo $(MAIN)
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-linux-$(VERSION).tgz $(BINARY)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo $(MAIN)
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-darwin-$(VERSION).tgz $(BINARY)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD)/$(BINARY).exe -ldflags $(LDFLAGS) -a -tags netgo $(MAIN)
	tar -C $(BUILD) -llzcvf $(DIST)/$(BINARY)-windows-$(VERSION).tgz $(BINARY).exe

.PHONY: bootstrap
bootstrap:
ifeq (,$(wildcard ./go.mod))
	go mod init dockerctl
endif
	go mod download

.PHONY: clean
clean:
	rm -rf _*
