GO:=go

.PHONY: test
test:
	$(GO) test -v ./...

.PHONY: generate
generate:
	$(GO) generate ./... -x

.PHONY: generateCheck
generateCheck:
	$(GO) generate -x ./... && git diff --exit-code
	code=$?
	git checkout -- .
	exit $(code)
