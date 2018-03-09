PROJECT_NAME = eve
PACKAGE = github.com/AAA-Intelligence/$(PROJECT_NAME)

all: build run

build:
	go build -o $(PROJECT_NAME) -v $(PACKAGE)

run:
	./$(PROJECT_NAME)

clean:
	go clean $(PACKAGE)

deps:
	go get "github.com/gorilla/websocket"
			