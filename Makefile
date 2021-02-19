SHELL=/bin/bash

install: build
	scripts/build/install

build:
	scripts/build/plugin

.PHONY: build install clean
