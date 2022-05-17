FROM golang:latest as BUILD
RUN apt-get update && \
    apt-get install -y xvfb wkhtmltopdf ghostscript
WORKDIR /

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
COPY local_config/config.yaml /etc/
RUN go build ./pkg/server/main.go

ENTRYPOINT ["./main"]