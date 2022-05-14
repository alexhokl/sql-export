build:
	go build ./...
install:
	go install
test:
	go test ./...
coverage:
	go test --cover ./...
coverage-html:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html
