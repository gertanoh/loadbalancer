# ################## Helpers ######################

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	
# ################## QA ######################
.PHONY: audit
audit:
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	# staticcheck ./...
	@echo 'Running tests'
	go test -race -vet=off ./...

	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify

.PHONY: test
test:
	@echo 'Running tests'
	go test -race -vet=off ./...

# ################## Build ######################
current_time = $(shell date --iso-8601=seconds)
git_desc = $(shell git describe --always --dirty --tags --long)
linker_flags = '-X main.buildTime=${current_time}  -X main.version=${git_desc}'
## build: build the loadbalancer
.PHONY: build-lb
build-lb:
	@echo 'Building lb'
	go build -ldflags=${linker_flags} -o=./lb loadbalancer/*.go

## build: build the backend servers
.PHONY: build-srv
build-srv:
	@echo 'Building servers'
	go build -ldflags=${linker_flags} -o=./be servers/*.go

.PHONY: build
build:
	go build -ldflags=${linker_flags} -o=./lb loadbalancer/*.go
	go build -ldflags=${linker_flags} -o=./be servers/*.go