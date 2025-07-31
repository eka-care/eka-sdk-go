SHELL := /bin/bash
.PHONY: govendor

govendor:
	go mod tidy -compat=1.24.4
	go mod vendor
	git add vendor