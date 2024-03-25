build:
	go build cmd/wrangler/wrangler.go

install:
	go install cmd/wrangler/wrangler.go

test:
	go test -v ./...
