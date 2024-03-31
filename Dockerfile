FROM golang:1.22-alpine AS builder
# --- GO BASE ---
RUN apk update && \
	apk add --no-cache \
	"git" \
    "ca-certificates"
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o qobuz-sync ./cmd && chmod o+x qobuz-sync

# --- RUNTIME ---
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/qobuz-sync /

ENTRYPOINT [ "/qobuz-sync"]
