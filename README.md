# Project EVE - lightweight version

**This version of eve does not require any authentication so be careful with exposing it to the internet.**

## Clone

```
git clone --single-branch -b lightweight-with-cors https://github.com/AAA-Intelligence/eve.git
```

## Build

IMPORTANT: copy model and chat-data into bot directory before building

`docker build -t eve-light .`

## Run

`docker run --name eve-light -p 3001:3001 -d eve-light`

Create a HTTP POST request to `http://localhost:3001/message-api` with following body:

```json
{
  "message":"Hello eve"
}
```
