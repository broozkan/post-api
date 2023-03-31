db-test:
	go clean -testcache && go test ./internal/couchbase -run TestCouchbase -v

unit-test:
	go clean --testcache && go test ./... -short

lint:
	golangci-lint run -c .dev/.golangci.yml

generate-mocks:
	mockgen -destination=internal/mocks/mock_post_service.go -package mocks broozkan/postapi/handlers PostServiceInterface
	mockgen -destination=internal/mocks/mock_post_repository.go -package mocks broozkan/postapi/internal/services RepositoryInterface
