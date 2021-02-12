
build-local: FORCE
	go build -o runnr-local ./cmd/runnr/runnr.go

reset-tmp: FORCE
	rm -rf tmp

test: FORCE
	go test ./...

FORCE: