UNAME := $(shell uname)
ifeq ($(UNAME), Linux)
	TARGET_PLATFORM := "linux"
else ifeq ($(UNAME), Darwin)
	TARGET_PLATFORM := "darwin"
else
	TARGET_PLATFORM := "unknown"
endif

build: .version
	env GOOS=linux GOARCH=amd64 go build -ldflags="-X 'khosef/pkg/about.version=$(shell cat .version)'" -o build/khosef-linux-amd64 ./cmd
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'khosef/pkg/about.version=$(shell cat .version)'" -o build/khosef-darwin-amd64 ./cmd
	env GOOS=linux GOARCH=amd64 go build -o build/kh-aws-linux-amd64 ./ext/aws/cmd
	env GOOS=darwin GOARCH=amd64 go build -o build/kh-aws-darwin-amd64 ./ext/aws/cmd

test:
	echo "tests are for the weak"

.version:
	echo "0.0+dev-$(shell date +%s)-build" > .version

clean:
	rm -rf ./build
	rm -f .version

install:
	cp build/khosef-$(TARGET_PLATFORM)-amd64 $(HOME)/.bin/kh
	cp build/kh-aws-$(TARGET_PLATFORM)-amd64 $(HOME)/.bin/kh-aws

default: build
