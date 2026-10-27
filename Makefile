BINARY = file-receiver

all:
	go build -o $(BINARY) .

install:
	mkdir -p $(HOME)/.local/bin
	go build -o $(HOME)/.local/bin/$(BINARY) .

run:
	go run .

clean:
	rm -f server $(BINARY)

.PHONY: all install run clean