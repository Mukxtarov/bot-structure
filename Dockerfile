FROM golang:latest

ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o bot ./cmd/main.go

CMD ["./bot"]