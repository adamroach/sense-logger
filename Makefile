SRC=$(shell find * -name '*.go')

all: run

run: logger
	./logger

logger: $(SRC)
	go build ./cmd/$@
