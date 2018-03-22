# Project EVE
Interactive Bot Chat

# How to run
## Installation

```
go get "github.com/AAA-Intelligence/eve"
```

## Run
After starting the webserver open http://localhost:8080/register in your browser to create an account.

### With Makefile
```
cd "%GOPATH%/src/github.com/AAA-Intelligence/eve"
make deps
make all
```
### Without Makefile
#### Windows
```
go build -o "eve.exe" "github.com/AAA-Intelligence/eve/cmd/eve" 
eve.exe -http 8080
```
#### macOS / linux
```
go build -o "eve" "github.com/AAA-Intelligence/eve/cmd/eve" 
./eve -http 8080
```

