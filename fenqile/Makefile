EXEC=$(shell basename ${PWD})

EXEC:main.go
	@go build -o $(EXEC)

test:$(EXEC)
	./$(EXEC) -f ~/input.html
