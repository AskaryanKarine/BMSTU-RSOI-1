FROM golang:1.22.5 AS builder

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY .. ./

RUN go build -o ./build/main ./cmd/person/main.go

FROM alpine:latest AS production

COPY --from=builder /app/build/main ./build/main

EXPOSE 8000
ENTRYPOINT ./build/main
