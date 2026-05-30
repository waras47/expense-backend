from golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/server .
COPY .env .

EXPOSE 8081

CMD ["./server"]