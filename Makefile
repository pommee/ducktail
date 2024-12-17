.PHONY: test

test:
	./testing/logs.sh | go run main.go
