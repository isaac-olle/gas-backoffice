FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /gas ./cmd/api/main.go

# ---

FROM scratch

COPY --from=builder /gas /gas
COPY --from=builder /app/config/config.yaml /config/config.yaml

ENV CONFIG_FILE_PATH=/config/config.yaml
ENV ENVIRONMENT=production

ENTRYPOINT ["/gas"]
