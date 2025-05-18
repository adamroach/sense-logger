SRC=$(shell find * -name '*.go')

all: run

run: logger
	./logger

logger: $(SRC)
	go build ./cmd/$@

format:
	find . -name '*.go' -exec go fmt '{}' \;
