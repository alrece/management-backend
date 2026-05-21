FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /job-service ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /job-service /usr/local/bin/
COPY configs/config.yaml /etc/job-service/config.yaml
EXPOSE 8082
ENTRYPOINT ["job-service"]
