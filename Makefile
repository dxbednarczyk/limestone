PREFIX = $(HOME)/.local

build:
	mkdir -p build
	go build -o build

install: build/limestone
	install -Dm755 build/limestone $(PREFIX)/bin/limestone

uninstall: $(PREFIX)/bin/limestone
	@echo -n "Do you want to remove all saved logins and configuration? [y/n] "
	@read line; if [ $$line = "y" ]; then rm -rf $$HOME/.config/limestone; fi

	rm -f $(PREFIX)/bin/limestone

clean:
	rm -rf build
	go mod tidy
	go clean
