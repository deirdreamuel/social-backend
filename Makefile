build:
	GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o main ./cmd/main.go