SHELL=/bin/bash

install: build
	scripts/build/install

build:
	scripts/build/plugin

release:
	scripts/build/release-build

.PHONY: build install clean
