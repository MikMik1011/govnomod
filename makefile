default: build

build:
	CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -buildmode=c-shared -o build/govnomod.so main.go

clean:
	rm -rf build