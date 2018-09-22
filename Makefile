PROJECT_NAME = eve
PACKAGE = github.com/AAA-Intelligence/$(PROJECT_NAME)
cmd = $(PACKAGE)/cmd/$(PROJECT_NAME)

ifeq ($(OS),Windows_NT)
    FILENAME = $(PROJECT_NAME).exe
else
    FILENAME = $(PROJECT_NAME)
endif

# run args
HTTP = 8080

all: build run

build:
	go build -o $(FILENAME) -v $(cmd)

run:
	./$(FILENAME) -http $(HTTP)

clean:
	go clean $(cmd)

deps:
	go get "github.com/gorilla/websocket" "golang.org/x/crypto/bcrypt" "github.com/drhodes/golorem" "github.com/gorilla/schema" "github.com/rs/cors"
			