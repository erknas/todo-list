build:
	@go build -o bin/todo-list cmd/main.go
	
run: build
	@./bin/todo-list