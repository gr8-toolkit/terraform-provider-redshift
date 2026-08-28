TEST?=$$(go list ./... | grep -v 'vendor')

.PHONY: build
build:
	go build -o ./dist/

.PHONY: test
test:
	go test ./redshift -v

.PHONY: testacc
testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -count=1 -timeout 120m

.PHONY: doc
doc:
	@go generate ./...
