.PHONY: build install test clean

INSTALL_DIR ?= /usr/local/bin

build:
	go build -o jig .

install: build
	install -m 755 jig $(INSTALL_DIR)/jig

install-local: build
	mkdir -p $(HOME)/.local/bin
	install -m 755 jig $(HOME)/.local/bin/jig

test:
	go test ./...

clean:
	rm -f jig
