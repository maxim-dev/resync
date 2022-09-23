help:
	@echo 'Targets:'
	@echo '  all          - download dependencies and compile binary'
	@echo '  deps         - download dependencies'
	@echo '  build        - compile binary'
	@echo '  test         - run short unit tests'
	@echo '  bench        - run benchmarks alongside tests'
	@echo '  tidy         - tidy go modules'
	@echo '  clean        - delete build artifacts'
	@echo ''
	@exit 0

all: deps build

deps:
	@echo ">> dependencies are downloading..."
	@go mod download -x
	@echo ">> done"

tidy:
	@echo ">> syncing dependencies..."
	@go mod verify
	@go mod tidy
	@echo ">> done"

generate:
	@go generate ./...

test: generate
	@go test ./...

run:
	@go run ./cmd/app/ $(filter-out $@,$(MAKECMDGOALS))

build:
	@echo ">> building..."
	@go build -o build/resync ./cmd/app/
	@echo ">> done. See ./build directory"

clean:
	@echo ">> cleaning..."
	@go clean
	@rm -rf ./build
	@echo ">> done"


bench: generate
	@go test ./... -bench . -benchmem #-benchtime=3s

