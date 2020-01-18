FROM golang:1.11 as builder
WORKDIR /root

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd cmd
COPY manager manager

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix nocgo -o eve ./cmd/eve


FROM python:3.6
WORKDIR /server

RUN pip install --upgrade pip

# Copy our project to the Go environment
COPY bot bot

# Install EVE dependencies
RUN pip install -r ./bot/requirements.txt

COPY --from=builder /root/eve .

# Define entrypoint for server application
ENTRYPOINT [ "./eve", "-http", "3001" ]
EXPOSE 3001
