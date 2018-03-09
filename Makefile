PROJECT_NAME = eve
PACKAGE = github.com/AAA-Intelligence/$(PROJECT_NAME)

# run args
HTTP = 8080

all: build run

build:
	go build -o $(PROJECT_NAME) -v $(PACKAGE)

run:
	./$(PROJECT_NAME) -http $(HTTP)

clean:
	go clean $(PACKAGE)

deps:
	go get "github.com/gorilla/websocket"
			