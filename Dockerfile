FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o panflow ./cmd/server

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/panflow .
COPY --from=builder /app/config.example.yaml ./config.example.yaml

EXPOSE 8080

ENTRYPOINT ["./panflow"]
