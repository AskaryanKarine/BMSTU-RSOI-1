FROM golang:1.22.5 AS builder

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY .. ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/main ./cmd/person/main.go

FROM alpine:latest AS app

COPY --from=builder /go/bin/main ./go/

EXPOSE 8000
CMD ./go/main
