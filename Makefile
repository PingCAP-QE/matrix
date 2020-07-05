GOARCH := $(if $(GOARCH),$(GOARCH),amd64)
GO=GO15VENDOREXPERIMENT="1" CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GO111MODULE=on go
GOTEST=GO15VENDOREXPERIMENT="1" CGO_ENABLED=1 GO111MODULE=on go test # go race detector requires cgo
VERSION   := $(if $(VERSION),$(VERSION),latest)
IMAGE := $(if $(IMAGE),$(IMAGE),pingcap/matrix)

GOBUILD=$(GO) build

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

PACKAGE_LIST := go list ./... | grep -vE "matrix/test"
PACKAGE_DIRECTORIES := $(PACKAGE_LIST) | sed 's|chaos-mesh/matrix/||'

default: all

all:
	$(GOBUILD) $(GOMOD) -o bin/matrix src/main.go

groupimports: $(GOBIN)/goimports
	$< -w -l -local chaos-mesh/matrix $$($(PACKAGE_DIRECTORIES))

fmt: groupimports
	go fmt ./...

tidy:
	@echo "go mod tidy"
	GO111MODULE=on go mod tidy
	@git diff --exit-code -- go.mod

docker-build:
	docker build -t $(IMAGE):$(VERSION) .

$(GOBIN)/goimports:
	$(GO) get golang.org/x/tools/cmd/goimports@v0.0.0-20200309202150-20ab64c0d93f

test:
	$(GOTEST) $$($(PACKAGE_LIST))
