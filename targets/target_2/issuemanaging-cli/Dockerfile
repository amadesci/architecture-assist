# issuemanaging_cli/Dockerfile
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /issuemanaging-cli

FROM alpine:latest
WORKDIR /app
COPY --from=build /issuemanaging-cli ./issuemanaging-cli
EXPOSE 8080
CMD ["./issuemanaging-cli"]