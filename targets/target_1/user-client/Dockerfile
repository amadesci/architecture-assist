# user_client/Dockerfile
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-client

FROM alpine:latest
WORKDIR /app
COPY --from=build /user-client ./user-client
EXPOSE 8080
CMD ["./user-client"]