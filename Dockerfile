FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o biblioteca-api ./cmd/api

FROM golang:1.26-alpine

WORKDIR /app

COPY --from=builder /app/biblioteca-api .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./biblioteca-api"]
