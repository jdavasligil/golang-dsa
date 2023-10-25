build:
		@go build -o ./target/golang-dsa ./cmd
run:
		@./target/golang-dsa
test:
		@go test -v ./...
