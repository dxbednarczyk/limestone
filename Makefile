PREFIX := $(HOME)/.local

build:
	mkdir -p build
	go build -o build

install: build/limestone
	install -Dm755 build/limestone $(PREFIX)/bin/limestone

uninstall: $(PREFIX)/bin/limestone
	rm -f $(PREFIX)/bin/limestone

clean: build/
	rm -rf build

.PHONY: reportcard
reportcard:
	gofmt -s -w -l .
	go vet

#   not necessarily goreportcard related, but important project cleaning tools
	go mod tidy
	go clean
