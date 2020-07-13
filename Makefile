BIN := infraview
TOOL_BIN := $(PWD)/bin

GOLINT := $(TOOL_BIN)/golint
GOIMPORTS := $(TOOL_BIN)/goimports
ENUMER := $(TOOL_BIN)/enumer

VERSION= $(shell git describe --tags --always)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X github.com/cycloidio/infraview/cmd.Version=${VERSION}"

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

.PHONY: test
test: ## Tests all the project
	@go test ./...

.PHONY: lint
lint: $(GOLINT) $(GOIMPORTS) ## Runs the linter
	@$(GOLINT) -set_exit_status ./... && test -z "`go list -f {{.Dir}} ./... | xargs $(GOIMPORTS) -l | tee /dev/stderr`"

.PHONY: generate
generate: $(ENUMER) ## Generates the needed code
	@go generate ./...

.PHONY: build
build: ## Builds the binary
	go build -o $(BIN) ${LDFLAGS}

.PHONY: install
install: ## Install the binary
	go install ${LDFLAGS}
