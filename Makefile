ERRCHECK_EXCLUDE_PATTERN := '(cmd/redirector-server|restapi/operations)'

dev-deps:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger
	go get -u github.com/golang/dep/cmd/dep
	go get -u golang.org/x/tools/cmd/goimports
	go get -u golang.org/x/lint/golint
	go get -u github.com/kisielk/errcheck

deps:
	dep ensure -vendor-only

gen-swagger:
	swagger generate server -f api.yml

lint:
	go vet ./...
	bash -c "test $$(goimports -d $$(git ls-files *.go) | tee /dev/stderr | wc -l) -eq 0"
	bash -c "test $$(go list ./... | xargs golint | tee /dev/stderr | wc -l) -eq 0"
	bash -c "errcheck \$$(go list ./... | grep -v -E $(ERRCHECK_EXCLUDE_PATTERN))"

run-test:
	go test ./...

clean:
	git clean -f -d -X  -- **
	git clean -f -d -X  -- cmd/**
