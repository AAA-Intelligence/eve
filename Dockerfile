FROM python:3.6
LABEL maintainer="Niklas Korz <nk@alugha.com>"
ENV GOPATH=/goenv

# Install general dependencies
RUN apt-get update
RUN apt-get install -y build-essential
RUN wget https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz
RUN tar -xvzf go1.10.1.linux-amd64.tar.gz -C /
ENV PATH="${PATH}:/go/bin"
RUN pip install --upgrade pip

# Copy our project to the Go environment
COPY . /goenv/src/github.com/AAA-Intelligence/eve
WORKDIR /goenv/src/github.com/AAA-Intelligence/eve

# Install EVE dependencies
RUN pip install -r ./bot/requirements.txt
RUN make deps

# Build the server
RUN make build

# Define entrypoint for server application
EXPOSE 8080
ENTRYPOINT [ "./eve", "-http", "8080" ]
