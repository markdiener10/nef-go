SHELL := /usr/bin/env bash

.DEFAULT_GOAL := test

.PHONY: clean test 

clean:
	go clean -testcache -cache -modcache
	go mod tidy	
	@echo "Project CLEANED ##############"

test: 
	go vet 
	go test -cover

