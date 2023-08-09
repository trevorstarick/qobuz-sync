install:
	go build -o $(shell go env GOBIN)/qobuz-sync cmd/main.go
