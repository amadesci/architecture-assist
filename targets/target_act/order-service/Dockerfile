# order-service/Dockerfile
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /order-service

FROM alpine:latest
WORKDIR /app
COPY --from=build /order-service ./order-service
EXPOSE 8080
CMD ["./order-service"]