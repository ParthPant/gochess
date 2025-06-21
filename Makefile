name := gochess
all: build run
build:
	@go build
test:
	@go test
run:
	@./$(name)
clean:
	@rm -f $(name)
