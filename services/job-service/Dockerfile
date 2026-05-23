FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY pkg/ /app/pkg/
COPY api/ /app/api/
COPY services/job-service/go.mod services/job-service/go.sum /app/services/job-service/

WORKDIR /app/services/job-service
RUN go mod download

COPY services/job-service/ /app/services/job-service/
COPY pkg/ /app/pkg/
COPY api/ /app/api/

RUN go mod tidy && CGO_ENABLED=0 go build -o /job-service ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /job-service /usr/local/bin/
COPY services/job-service/configs/config.yaml /etc/job-service/config.yaml
EXPOSE 8082
ENTRYPOINT ["job-service"]
