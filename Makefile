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

.PHONY: lint
lint:
	gci write --skip-generated . 
	gofumpt -l -w .
	
	go vet
	go mod tidy
	go clean
