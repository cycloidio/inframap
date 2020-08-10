BIN := inframap
TOOL_BIN := $(PWD)/bin

GOLINT := $(TOOL_BIN)/golint
GOIMPORTS := $(TOOL_BIN)/goimports
ENUMER := $(TOOL_BIN)/enumer
PKGER := $(TOOL_BIN)/go-bindata

ARCHITECTURES=386 amd64
PLATFORMS=darwin linux windows

VERSION= $(shell git describe --tags --always)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X github.com/cycloidio/inframap/cmd.Version=${VERSION}"

.PHONY: help
help: Makefile ## This help dialog
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
	for help_line in $${help_lines[@]}; do \
		IFS=$$'#' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf "%-30s %s\n" $$help_command $$help_info ; \
	done

$(ENUMER):
	@GOBIN=$(TOOL_BIN) go install github.com/dmarkham/enumer

$(GOLINT):
	@GOBIN=$(TOOL_BIN) go install golang.org/x/lint/golint

$(GOIMPORTS):
	@GOBIN=$(TOOL_BIN) go install golang.org/x/tools/cmd/goimports

$(PKGER):
	@GOBIN=$(TOOL_BIN) go install github.com/markbates/pkger/cmd/pkger

.PHONY: test
test: ## Tests all the project
	@go test ./...

.PHONY: lint
lint: $(GOLINT) $(GOIMPORTS) ## Runs the linter
	@$(GOLINT) -set_exit_status ./... && test -z "`go list -f {{.Dir}} ./... | xargs $(GOIMPORTS) -l | tee /dev/stderr`"

.PHONY: generate
generate: $(ENUMER) ## Generates the needed code
	@go generate ./...

.PHONY: generate-icons
generate-icons: $(PKGER) ## Generates the needed code and Icons, it's separated as the icons generate always a new output
	@go generate -tags icons ./... 

.PHONY: build
build: ## Builds the binary
	go build -o $(BIN) ${LDFLAGS}

.PHONY: install
install: ## Install the binary
	go install ${LDFLAGS}

.PHONY: build-all build-compress
build-all: ## Builds the binaries
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES),\
	$(shell export GO111MODULE=on; export CGO_ENABLED=0; export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o $(BUILD_PATH)/$(BIN)-$(GOOS)-$(GOARCH) ${LDFLAGS})))

build-compress: build-all ## Builds and compress the binaries
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES),\
	$(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); tar -C $(BUILD_PATH) -czf $(BUILD_PATH)/$(BIN)-$(GOOS)-$(GOARCH).tar.gz $(BIN)-$(GOOS)-$(GOARCH))))

.PHONY: dbuild
dbuild: ## Builds the docker image with same name as the binary
	@docker build -t $(BIN) .
