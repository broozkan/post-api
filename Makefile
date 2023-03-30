unit-test:
	go clean --testcache && go test ./... -short

lint:
	golangci-lint run -c .dev/.golangci.yml

generate-mocks:
	mockgen -destination=internal/mocks/mock_handler_service.go -package mocks broozkan/redditpostapi/internal/handlers PostServiceInterface
