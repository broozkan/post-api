code-coverage:
	go test `go list ./... | grep -vE "/tilt_modules|/metrics|/contract|/mocks|/repository"` -coverprofile cover.out
	go tool cover -html=cover.out -o coverage.html
	echo `go tool cover -func cover.out | grep total`

db-test:
	go clean -testcache && go test ./internal/couchbase -run TestCouchbase -v

unit-test:
	go clean --testcache && go test ./... -short

lint:
	golangci-lint run -c .dev/.golangci.yml

generate-mocks:
	mockgen -destination=internal/mocks/mock_post_service.go -package mocks broozkan/postapi/handlers PostServiceInterface
	mockgen -destination=internal/mocks/mock_post_repository.go -package mocks broozkan/postapi/internal/services RepositoryInterface
