# Builder stage
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server

# Final stage
FROM gcr.io/distroless/static-debian11 AS build-release-stage
# If HTTPS listener needs CA certs, uncomment:
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/server /server

# Expose HTTPS (443) and DNS (53 TCP/UDP)
# EXPOSE 443
# EXPOSE 53/tcp 53/udp
EXPOSE 80
EXPOSE 52/tcp 52/udp

ENTRYPOINT ["/server"]
