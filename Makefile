TEST_RESULTS=/tmp/test-results

fmt:
	go fmt ./...

test:
	go test ./...

build:
	go build ./...

mock:
	mockgen -source counter/counter.go -destination counter/mock_counter/mock_counter.go -package mock_counter

coverage:
	@mkdir -p ${TEST_RESULTS}
	@go test ./... -coverprofile=${TEST_RESULTS}/unittest.out -v $(GOPACKAGES)
	@go tool cover -html=${TEST_RESULTS}/unittest.out -o ${TEST_RESULTS}/unittest-coverage.html
	@rm -f ${TEST_RESULTS}/unittest.out