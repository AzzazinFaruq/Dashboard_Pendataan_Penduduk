# --- Build stage ---
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# --- Runtime stage ---
FROM alpine:latest
WORKDIR /app
# CA certs untuk koneksi TLS bila dibutuhkan
RUN apk --no-cache add ca-certificates
COPY --from=build /app/app .
COPY --from=build /app/public ./public
EXPOSE 8080
CMD ["./app"]
