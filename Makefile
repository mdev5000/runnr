
build-local: FORCE
	go build -o runnr-local ./runnr/runnr.go

reset-tmp: FORCE
	rm -rf tmp

FORCE: