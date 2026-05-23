FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY pkg/ /app/pkg/
COPY api/ /app/api/
COPY services/workflow-service/go.mod services/workflow-service/go.sum /app/services/workflow-service/

WORKDIR /app/services/workflow-service
RUN go mod download

COPY services/workflow-service/ /app/services/workflow-service/
COPY pkg/ /app/pkg/
COPY api/ /app/api/

RUN go mod tidy && CGO_ENABLED=0 go build -o /workflow-service ./cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /workflow-service /usr/local/bin/
COPY services/workflow-service/configs/config.yaml /etc/workflow-service/config.yaml
EXPOSE 8083
ENTRYPOINT ["workflow-service"]
