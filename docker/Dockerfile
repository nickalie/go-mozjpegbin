FROM golang

RUN apt-get update && apt-get install -y --no-install-recommends libjpeg-turbo-progs && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go test -v ./...
