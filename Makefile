.PHONY: all
all: goat

goat: goat.go
	CGO_ENABLED=0 go build -o $@ goat.go

.PHONY: clean
clean:
	go clean

