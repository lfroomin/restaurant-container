NAME=restaurant-container
VERSION=0.0.1

.PHONY: build
## build: Compile the packages.
build:
	@go build -o $(NAME)

.PHONY: run
## run: Build and Run.
run: build
	@./$(NAME)

.PHONY: clean
## clean: Clean project and previous builds.
clean:
	@rm -f $(NAME)

.PHONY: deps
## deps: Download modules
deps:
	@go mod download

.PHONY: test
## test: Run tests with verbose mode
test:
	@go test ./...

.PHONY: docker-build
## test: Build the Docker image
docker-build:
	@docker build --tag restaurant-container:alpha .

.PHONY: docker-run
## test: Run the Docker container
docker-run:
	@docker run --rm -p 8080:8080 --name restaurant-container -e AWS_REGION=us-west-2 \
     -v ~/.aws/credentials:/root/.aws/credentials:ro restaurant-container:alpha


.PHONY: help
all: help
# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
