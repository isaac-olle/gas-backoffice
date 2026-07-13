FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /gas ./cmd/api/main.go

# ---

FROM scratch

COPY --from=builder /gas /gas
COPY --from=builder /app/config/config.yaml /config/config.yaml
COPY --from=builder /app/deploy/golang-migrate /golang-migrate

ENV CONFIG_FILE_PATH=/config/config.yaml
ENV MIGRATIONS_PATH=/golang-migrate

ENTRYPOINT ["/gas"]
