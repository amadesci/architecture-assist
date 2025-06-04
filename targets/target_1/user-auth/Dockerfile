# user_auth/Dockerfile
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-auth

FROM alpine:latest
WORKDIR /app
COPY --from=build /user-auth ./user-auth
EXPOSE 8080
CMD ["./user-auth"]