# Project EVE - lightweight version

**This version of eve does not require any authentication so be careful with exposing it to the internet.

## Build

IMPORTANT: copy model and chat-data into bot directory before building

`docker build -t eve-light .`

## Run

`docker run --name eve-light -p 8080:8080 -d eve-light`

Create a HTTP POST request to http://localhost:8080/messageApi with following body:

```json
{
  "message":"Hello eve"
}
```
